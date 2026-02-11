# Introduction

a simple REST API written in `Golang` that exposes two primary endpoints: one for
creating a new user, and the other for logging in using username and password.

`PostgresDB` is used as the database of choice.

`ELK` stack is used for log management, while `Prometheus` and `Grafana` are used for analysing metrics.

The aim is to iteratively scale this simple login-endpoint
to handle a large number of requests, learning the process involved while doing so.

[tsenart/vegeta](https://github.com/tsenart/vegeta) will be used for load testing and measuring the performance.