package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	BatchInsertSeconds int      `json:"batchInsertSeconds"`
	Port               int      `json:"port"`
	DbConfig           DbConfig `json:"database"`
}

var configFilePath string

type DbConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type HrefClick struct {
	IpAddress  string `json:"ipAddress"`
	Url        string `json:"url"`
	Href       string `json:"href"`
	HrefTop    int    `json:"hrefTop"`
	HrefRight  int    `json:"hrefRight"`
	HrefBottom int    `json:"hrefBottom"`
	HrefLeft   int    `json:"hrefLeft"`
}

var hrefClicks []HrefClick

type PageView struct {
	Domain       string `json:"domain"`
	IpAddress    string `json:"ipAddress"`
	Url          string `json:"url"`
	UserAgent    string `json:"userAgent"`
	ScreenHeight int    `json:"screenHeight"`
	ScreenWidth  int    `json:"screenWidth"`
}

var pageViews []PageView

func listenForRecords(db *sql.DB, seconds time.Duration) {
	// Run every x seconds.
	for _ = range time.Tick(seconds) {
		// Handle page views.
		newPageViews := make([]PageView, len(pageViews))
		copy(newPageViews, pageViews)
		go SetPageViews(db, newPageViews)
		pageViews = pageViews[0:0]

		// Handle href clicks.
		newHrefClicks := make([]HrefClick, len(hrefClicks))
		copy(newHrefClicks, hrefClicks)
		go SetHrefClicks(db, newHrefClicks)
		hrefClicks = hrefClicks[0:0]
	}
}

func IpAddress(remoteAddr string) string {
	arr := strings.Split(remoteAddr, ":")
	return arr[0]
}

func hrefClickHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	hrefClick := HrefClick{}
	if err := json.Unmarshal(body, &hrefClick); err != nil {
		log.Println("Unable to unmarshal hrefClick: ", err)
	}
	// Get ip address from http request
	hrefClick.IpAddress = IpAddress(r.RemoteAddr)
	hrefClicks = append(hrefClicks, hrefClick)
	w.WriteHeader(201)
}

func pageViewsHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	pageView := PageView{}
	if err := json.Unmarshal(body, &pageView); err != nil {
		log.Println("Unable to unmarshal pageView: ", err)
	}
	// Get ip address from http request
	pageView.IpAddress = IpAddress(r.RemoteAddr)
	pageViews = append(pageViews, pageView)
	w.WriteHeader(201)
}

func readConfig(configFilePath string) Config {
	config := Config{}
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal("Unable to read config file: ", err)
	}
	if err = json.Unmarshal(configFile, &config); err != nil {
		log.Fatal("Unable to unmarshal configFile into config: ", err)
	}
	return config
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, []byte)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with, x-requested-by, Content-Type")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Unable to read requeset body: ", err)
		}
		fn(w, r, body)
	}
}

func init() {
	goPath := os.Getenv("GOPATH")
	defaultConfigPath := fmt.Sprintf("%s/src/github.com/roberttstephens/webanalytics/config.json", goPath)
	flag.StringVar(&configFilePath, "config", defaultConfigPath, "path to config.json")
}

func main() {
	// Read the config, initialize the database and listen for records.
	flag.Parse()
	config := readConfig(configFilePath)
	db := Db(config.DbConfig)
	seconds := time.Duration(config.BatchInsertSeconds) * time.Second
	go listenForRecords(db, seconds)

	// Create the handlers for page-view/ and href-click/ POSTs
	http.HandleFunc("/page-views/", makeHandler(pageViewsHandler))
	http.HandleFunc("/href-click/", makeHandler(hrefClickHandler))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
