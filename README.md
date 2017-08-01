[![Build Status](https://travis-ci.org/jeremy-miller/weather.svg?branch=master)](https://travis-ci.org/jeremy-miller/weather)
[![Coverage Status](https://coveralls.io/repos/github/jeremy-miller/weather/badge.svg?branch=master)](https://coveralls.io/github/jeremy-miller/weather?branch=master)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/jeremy-miller/weather/blob/master/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.8-blue.svg)]()

# Weather
API server for retrieving the weather of a particular city.
This project is based on [this](http://howistart.org/posts/go/1/) tutorial.

<details>
<summary><strong>Table of Contents</strong></summary>

* [Motivation](#motivation)
* [Usage](#usage)
  + [Prerequisites](#prerequisites)
  + [Build](#build)
  + [Static Code Analysis](#static-code-analysis)
  + [Test](#test)
  + [Run](#run)
  + [Example Call](#example-call)
* [License](#license)
</details>

## Motivation
I created this project to try Go.

## Usage
This implementation uses a Docker container to isolate the execution environment.

### Prerequisites
- [Docker](https://docs.docker.com/engine/installation/)

### Build
Before interacting with the weather server, the Docker container must be built: ```docker build -t jeremymiller/weather .```

### Static Code Analysis


### Test
To run the Weather tests, execute the following command: ```docker run -it --rm jeremymiller/weather go test```

### Run
1. Start the weather server by executing the following command: ```docker run -it --rm -p 8080:8080 jeremymiller/weather```
2. Call the REST API for the city of your choice to `http://localhost:8080/weather/<city>`

### Example Call
```
$ curl http://localhost:8080/weather/tokyo
{"city":"tokyo","temp":304.3,"took":"666.678077ms"}
```
*NOTE: The temperature is returned in Kelvin*

## License
[MIT](https://github.com/jeremy-miller/weather/blob/master/LICENSE)
