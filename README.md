# Notion AI

- [Description](#description)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
  - [HTTP](#http)
    - [Start HTTP Server](#start-http-server)
    - [Request Prompt](#request-prompt)

## Description

This application integrates Notion AI into your workflow. It has three main functionalities:

![In-use Animation](https://github.com/codercola034/notion-ai/blob/main/demo.gif?raw=true "In-use Animation")

## Features

- **HTTP Server**: It can run as an HTTP server, receiving and responding to requests.
- **Terminal User Interface (TUI)**: It allows user interaction in a terminal-based interface.
- **Normal Prompt**: It can run as a normal prompt for command-line inputs.

## Getting Started

### Installation

```shell
go install github.com/codercola034/notion-ai@latest
```

### Usage

```shell
Usage of notion-ai
  -http
        run in http mode
  -port int
        port to run http server on (default 8080)
  -tui
        run in terminal user interface mode
  -prompt string
        run with custom prompt
```

### HTTP

#### Start Http Server

```shell
notion-ai -http -port 10000
```

#### Request Prompt

There are only one endpoint listening /prompt

```shell
curl http://localhost:10000/prompt -d "write a golang hello world programme"

```

_Response_

````shell
Here is a simple "Hello World" program written in Golang:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

```

This program prints "Hello, World!" to the console.

````
