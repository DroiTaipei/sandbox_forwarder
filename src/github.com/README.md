# Introduction
The configurations, constants, error code, restful API payloads, and shared utilities of Golang based microservices in Droi. Require `Golang v1.8.0` or above.

## Import the package to your project
- Instruction on [KnowledgeBase](http://mkdocs:8000/development/Golang/create_new_project/)

## Update submodule
```
# In project folder
$ git pull
$ git submodule foreach git pull
```

# Project Structure
```
├── CHANGELOG.md
├── README.md
├── cloudcache
│   ├── app.go
│   └── app_test.go
├── errors.go
├── errors_test.go
├── file
│   ├── cdnpublisher
│   │   └── struct.go
│   ├── error.go
│   ├── majesty
│   │   └── struct.go
│   ├── utils.go
│   ├── utils_test.go
│   └── variable.go
├── logger.go
└── util
    └── slack
```

# Development
The package introduces [testify](https://github.com/stretchr/testify/assert) and [httpmock](https://github.com/jarcoal/httpmock) for unit testing, and you can import them by `go get`
```
$ go get github.com/stretchr/testify/assert github.com/jarcoal/httpmock
```
## The concept about error

### DroiError
* The basic error type (interface)
* Two Elements
    * Error() - error message
    * ErrorCode() - error code

### AsDroiError
* According to the name, you will know it could be handle as Droi Error
* But it was created to support more actually.
* The Implemented struct type is TraceDroiError.
* Interface Detail
    * AsEqual - It should contain a DroiError, AsEqual should be used in error handling.