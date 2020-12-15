package watcher

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

type Values struct {
	Ping    time.Duration
	DLSpeed float64
	ULSpeed float64
}

type Watcher struct {
	serverCount     int
	serverBlacklist []string

	LastValues *Values
	ticker     *time.Ticker

	ErrorCount uint64
	OnError    func(err error, data interface{})
}

func NewWatcher(serverCount int, serverBlacklist []string, fetchInterval time.Duration) *Watcher {
	for i, bl := range serverBlacklist {
		serverBlacklist[i] = strings.TrimSpace(strings.ToLower(bl))
	}

	ticker := time.NewTicker(fetchInterval)

	w := &Watcher{
		serverCount,
		serverBlacklist,
		new(Values),
		ticker,
		0,
		func(error, interface{}) {},
	}

	go w.watchWorker()

	return w
}

func (w *Watcher) FetchValuesBlocking() (err error) {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return
	}

	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		return
	}

	if w.serverCount == 0 {
		w.serverCount = len(serverList.Servers)
	}

	count := 0
	var pings time.Duration
	var dlSpeeds float64
	var ulSpeeds float64

	for _, server := range serverList.Servers {
		if w.isBlacklisted(server) {
			continue
		}

		if count >= w.serverCount {
			break
		}

		if err = testServer(server); err != nil {
			w.onErrorWrapper(err, server)
			continue
		}

		pings += server.Latency
		dlSpeeds += server.DLSpeed
		ulSpeeds += server.ULSpeed

		count++
	}

	w.LastValues.Ping = time.Duration(pings.Nanoseconds()/int64(count)) * time.Nanosecond
	w.LastValues.DLSpeed = dlSpeeds / float64(count)
	w.LastValues.ULSpeed = ulSpeeds / float64(count)

	return
}

func (w *Watcher) onErrorWrapper(err error, v interface{}) {
	atomic.AddUint64(&w.ErrorCount, 1)
	if w.OnError != nil {
		w.OnError(err, v)
	}
}

func (w *Watcher) watchJob() {
	if err := w.FetchValuesBlocking(); err != nil {
		w.onErrorWrapper(err, nil)
	}
}

func (w *Watcher) watchWorker() {
	for {
		go w.watchJob()
		<-w.ticker.C
	}
}

func (w *Watcher) isBlacklisted(server *speedtest.Server) bool {
	for _, bl := range w.serverBlacklist {
		if strings.ToLower(server.Name) == bl || server.Host == bl {
			return true
		}
	}
	return false
}

func testServer(server *speedtest.Server) (err error) {
	if err = server.PingTest(); err != nil {
		return
	}

	if err = server.DownloadTest(); err != nil {
		return
	}

	err = server.UploadTest()

	return
}
