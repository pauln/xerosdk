# XERO SDK (alpha)

The aim of this repository is to build a new SDK for Xero platform based on the OAuth2 protocol.

We will take some parts of the old repository for Xero Go SDK https://github.com/XeroAPI/xerogolang

The work is still in progress

Documentation reference https://developer.xero.com/documentation/


### Configure

The Xero Golang SDK can be configured in multiple ways using the Provider type what accepts a Config type with the minimum
information required to setup a new OAuth2 client.

If you want to run the /example you must use a .env with the next vars

```
CLIENT_ID="----"
CLIENT_SECRET="----"
SCOPES="-----" // Comma separated fields
REDIRECT_URL="-------"
```

### Example App

This repo includes an Example App that shows you how to use this SDK. The app contains example of most of the functions
available for the API.

To run the example app do the following:
```text
$ cd example
$ go run *.go
```

Now open up your browser and go to [http://localhost:3000](http://localhost:3000) to see the example.

### How to Contribute

Make a pull request...
