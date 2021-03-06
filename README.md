#webanalytics

##Overview
If you want web analytics, you can use one or more of several third party services. Webanalytics is a simple performant open source application that covers some common use cases.

##Use cases

- How many page views am I getting? (Sometimes it's difficult to tell with varnish)
- On which URIs?
- What percentage of users are still on IE x?
- Which content do users click on?

##About the project

- Uses go (golang) to process requests, via a RESTful API.
- Uses a postgresql database. The database design is purposefully simple in order to be efficient with writes.
- Uses javascript to submit posts.

##How to use

Webanalytics is broken into two parts. The server side application and the javascript.

###Server side application
If you have already set up your $GOPATH and added $GOPATH/bin to your $PATH you should:
- Create a postgres user and database for webanalytics.
- run "go get github.com/roberttstephens/webanalytics" without quotation marks.
- Copy $GOPATH/src/github.com/roberttstephens/webanalytics/config.json to somewhere of your choice.
- Edit config.json to reflect your new database connection and desired port.
- Run "webanalytics --config path/to/config.json" without quotation marks.


###Javascript
The javascript is in poor shape right now. However, you should be able to copy docs/webanalytics.js to your site, change your domain (and possibly port) and start receiving POSTs.  Please reach out to me if something doesn't work, so I can fix it.
