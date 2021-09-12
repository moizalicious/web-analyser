# Web Analyser
Web analyser is a basic web application that can be used to crawl and identify information of a web page from a given URL.

It is written using Go and the [Gin Web Framework](https://github.com/gin-gonic/gin).

# Building From Source (Linux)
Pre-requisites:
* Git
* Go 1.16 or higher

Steps:
```bash
# Clone the git repository.
~$ git clone https://github.com/moizalicious/web-analyser

# Move to root directory of the project.
~$ cd web-analyser

# Build the application.
~$ go build -o application
```

# Run Build
To simply run the web application, execute the generated binary after building:
```
~$ ./application
```

Once the application starts running, open the browser and go to `http://localhost:8080` to access the web view.

# Release Mode & Custom Ports
By default, the application will start running in `DEBUG` mode with port `8080`. This can be configured via specific environment variables.

To define a custom application port, an `APP_PORT` environment variable must be set with the custom application port before running:
```bash
~$ export APP_PORT=3000
```

To run the application in `RELEASE` mode, an `APP_MODE` environment variable must be set with the variable set as `release` before running:
```bash
~$ export APP_MODE=release
```

# Screenshots
