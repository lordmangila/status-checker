# Go Site Health Checker

This is a demo website health checker application utilising
[gorilla/websocket](https://github.com/gorilla/websocket) package.

The applications `Client` and `Server` pattern was greatly influenced from
the [gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/master/examples/chat).

The application also follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Usage

#### Installing dependencies

```bash
    $ go get github.com/gorilla/websocket
```

#### Checkout and running the application

```bash
    $ go get github.com/lordmangila/status-checker/cmd/checker
    $ cd `go list -f '{{.Dir}}' github.com/lordmangila/status-checker/cmd/checker` && cd ../..
    $ go run cmd/checker/*.go
```

> Open [http://localhost:8080/](http://localhost:8080/) in your browser.

## Overview

> The websocket server is available via [ws://localhost:8080/check](ws://localhost:8080/check).

> The websocket server accepts a valid `url` as a message.
> A validation check is applied to determine the `url`'s validity.
> Once valid, the server checks the `url` status every 5 minutes and updates the client until the connection is closed.

#### Sample valid url response
```json
    {
        "URL": "http://www.google.com",
        "StatusCode": 200,
        "Active": true,
        "Valid": true,
        "Error": ""
    }
```

#### Sample invalid url response
```json
    {
        "URL": "invalidurl",
        "StatusCode": 0,
        "Active": false,
        "Valid": false,
        "Error": "Invalid URI: invalidurl"
    }
```
