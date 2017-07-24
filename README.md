[![Code Climate](https://codeclimate.com/github/jeremy-miller/weather/badges/gpa.svg)](https://codeclimate.com/github/jeremy-miller/weather)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/jeremy-miller/weather/blob/master/LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.8.1-blue.svg)]()

# Weather
API server for retrieving the weather of a particular city.
This project is based on [this](http://howistart.org/posts/go/1/) tutorial.

- [Motivation](#motivation)
- [Usage](#usage)
  - [Prerequisites](#prerequisites)
  - [Build](#build)
  - [Run](#run)
  - [Example Call](#example-call)
- [License](#license)

## Motivation
I created this project to learn Go.

## Usage
This implementation uses a Docker container to isolate the execution environment.

### Prerequisites
- [Docker](https://docs.docker.com/engine/installation/)

### Build
Before interacting with the weather server, the Docker container must be built: ```docker build -t jeremymiller/weather .```

### Run
1. Start the weather server by executing the following command: ```docker run -it --rm jeremymiller/weather```
2. Call the REST API for the city of your choice to `http://localhost:8080/weather/<city>`

### Example Call
```
$ curl http://localhost:8080/weather/tokyo
{"city":"tokyo","temp":304.3,"took":"666.678077ms"}
```
*NOTE: The temperature is returned in Kelvin*

## License
[MIT](https://github.com/jeremy-miller/portals/blob/master/LICENSE)
