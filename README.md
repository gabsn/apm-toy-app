# apm-toy-app
A simple app to trace.

## Prerequisites
1. Install `docker` and `docker-compose`
2. (Optional) Install `go` and its [dependency management tool](https://github.com/golang/dep)

## Usage
1. Run the app
```
docker-compose up
```
This command should launch a minimal golang server along with a redis and a postgres instance.

2. Go to http://localhost:8080/, you should get something like that:
```
(247 hits) - City: Utrecht, 234323 inhabitants
```
Each time you hit this URL, the golang server will return a different city and its population from postgres and will also tell you how many times you hit this endpoint.

## Possible issues
If you encounter some dependency issues, try to run `dep ensure` (make sure you have it installed).
