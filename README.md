<div align="center">
    <h1>~ netrics ~</h1>
    <strong>
        speedtest.net + prometheus = netrics
    </strong>
<br>
</div>

---

# What is netrics?

netrics is a tool which allows to monitor your internet network speed using the [speedtest.net API](github.com/showwin/speedtest-go) in between regular time delays. On every test run, the average downstream, upstream and RTT is tested over a configurable ammount of servers tested subsequently.

netrics exposes a [Prometheus](https://prometheus.io/) metrics endpoint which makes the measured data accessable to be scraped by Prometheus.

The collected metrics can then be visualized using something like Grafana, for example.

![](https://i.imgur.com/b3j7rY4.png)

# Setup

The best use-case setup is to install [Docker](https://www.docker.com/) and [docker-compose](https://github.com/docker/compose) and then spin up the provided [`docker-compose.yml`](docker-compose.yml) which includes netrics, Prometheus and Grafana. 
```
$ docker-compose up -d
```

You can also use the provided [`traefik.docker-compose.yml`](traefik.docker-compose.yml) configuration which includes traefik as edge router and an appropriate configuration to expose all services over a single entrypoint.
```
$ docker-compose -f traefik.docker-compose.yml up -d
```

# Configuration

netrics is simply configured via the following command line flags:

```
Usage of netrics:
  -addr string
        Bind address of the mectric HTTP server. (default ":9091")
  -blacklist string
        Comma seperated list of excluded server naes or host addresses.
  -endpoint string
        HTTP (default "/metrics")
  -interval string
        Time interval between server fetches. (default "30m")
  -sc int
        Count of servers used from the server list to average values over. (default 3)
```

---

© 2020 Ringo Hoffmann (zekro Development)  
Covered by the MIT License.