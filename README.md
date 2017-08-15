# apm-toy-app
A simple app to trace.

![full_trace](https://github.com/gabsn/apm-toy-app/blob/img/full_trace.png)

## Prerequisites
1. Install `docker` and `docker-compose` (you can try `sudo pip install docker-compose`)
2. (Optional) Install `go` and its [dependency management tool](https://github.com/golang/dep)
3. You need a [Datadog account](https://www.datadoghq.com). If you don't have one, just use the free trial. 

## Get the app running
1. Get your Datadog API key (*Integrations > APIs* in the app) and paste it [here](https://github.com/gabsn/apm-toy-app/blob/master/docker-compose.yml#L28)

*Note:* each time you'll change your API key, you'll need to rebuild the datadog-agent docker image.

1. Get your Datadog API key (*Integrations > APIs*) and paste it [here](https://github.com/gabsn/apm-toy-app/blob/master/docker-compose.yml#L28)
2. Run the app
```
docker-compose up
```
This command should launch a minimal golang server along with redis, postgres and the datadog-agent.

3. Go to http://localhost:8080/, you should get something like that:
```
(247 hits) - City: Utrecht, 234323 inhabitants
```
Each time you hit this URL, the golang server will return a different city and its population from postgres and will also tell you how many times you hit this endpoint.

## Check the datadog-agent status
1. Connect to the container running the datadog-agent
```
docker exec -it apmtoyapp_datadog_1 bash
```
2. Run the info command of the agent
```
service datadog-agent info
```
3. Check that postgres and redis are properly reporting metrics
```
 Checks
  ======

    postgres (5.16.0)
    -----------------
      - instance #0 [OK]
      - Collected 15 metrics, 0 events & 1 service check

    redisdb (5.16.0)
    ----------------
      - instance #0 [OK]
      - Collected 36 metrics, 0 events & 1 service check
      - Dependencies:
          - redis: 2.10.5
```
4. Go in the Datadog app, *Metrics > Explorer* and check you can see your metrics in the app.

![metrics explorer](https://github.com/gabsn/apm-toy-app/blob/img/metrics_explorer.png)

## Possible issues
- If you encounter some dependency issues, try to run `dep ensure` (make sure you have [dep](https://github.com/golang/dep) installed).

- If you want to make modifications to the datadog image, you have to manually rebuild it with `docker build --no-cache datadog -t apmtoyapp_datadog`.
