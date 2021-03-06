Crypto Trading
[![Build Status](https://travis-ci.com/fabioDMFerreira/crypto-trading.svg?branch=master)](https://travis-ci.com/fabioDMFerreira/crypto-trading)
[![codecov](https://codecov.io/gh/fabioDMFerreira/crypto-trading/branch/master/graph/badge.svg)](https://codecov.io/gh/fabioDMFerreira/crypto-trading)
[![Maintainability](https://api.codeclimate.com/v1/badges/c3918bf55a339713d017/maintainability)](https://codeclimate.com/github/fabioDMFerreira/crypto-trading/maintainability)
====

Use smart rules to buy and to sell crypto assets.

## Features

* Connects with broker to buy/sell tokens (currently it only supports Kraken);
* Sends automatic events reports via email;
* Benchmarks algorithm.

## Technologies

* Go
* MongoDB
* React

## Executables

* `serviced` listens prices changes in broker, buys and sells based on algorithm rules, and sends events report.
* `benchmark` executes trading algorithm with multiple input combinations and publish the outputs into a CSV file.
* `webserver` starts an HTTP server with an API to get benchmark results and to execute benchmarks.
* `get-asset-prices` get prices from an external source and store in CSV files.
* `save-asset-prices` use CSV files to store prices in database.

### Setup serviced

Run application
```
$ go run cmd/serviced/main.go
```

### Setup webserver

Install dependencies
```
$ cd client && npm install
```

Run application
```
$ docker-compose up mongo
$ go run cmd/webserver/main.go
$ cd client && npm start
```
