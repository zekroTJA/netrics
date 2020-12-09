package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/showwin/speedtest-go/speedtest"
	"github.com/zekroTJA/netrics/internal/watcher"
)

var (
	flagAddr        = flag.String("addr", ":9091", "Bind address of the mectric HTTP server.")
	flagEndpoint    = flag.String("endpoint", "/metrics", "HTTP")
	flagServerCount = flag.Int("sc", 3, "Count of servers used from the server list to average values over.")
	flagBlacklist   = flag.String("blacklist", "", "Comma seperated list of excluded server naes or host addresses.")
	flagInterval    = flag.String("interval", "30m", "Time interval between server fetches.")
)

func main() {
	flag.Parse()

	interval, err := time.ParseDuration(*flagInterval)
	if err != nil {
		log.Fatalf("Argument Error: %s", err.Error())
	}

	blacklist := strings.Split(*flagBlacklist, ",")

	w := watcher.NewWatcher(*flagServerCount, blacklist, interval)

	w.OnError = func(err error, data interface{}) {
		switch d := data.(type) {
		case *speedtest.Server:
			log.Printf("[ERR] Server Error {%s, %s}: %s",
				d.Name, d.Host, err.Error())
		default:
			log.Printf("[ERR] %s", err.Error())
		}
	}

	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "netrics_average_rtt",
			Help: "Average server RTT in milliseconds.",
		}, w.RTTHandler),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "netrics_average_download_speed",
			Help: "Average server downlaod speed in MBit/s.",
		}, w.DLSpeedHandler),

		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "netrics_average_upload_speed",
			Help: "Average server upload speed in MBit/s.",
		}, w.ULSpeedHandler),

		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Name: "netrics_error_count_total",
			Help: "Total error count.",
		}, w.ErrorCountHandler),
	)

	mux := http.NewServeMux()
	mux.Handle(*flagEndpoint, promhttp.Handler())

	log.Printf("Metrics server listening on %s...", *flagAddr)
	if err = http.ListenAndServe(*flagAddr, mux); err != nil {
		log.Fatalf("Fatal: %s", err.Error())
	}
}
