// by liudan
package common

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"third/g2s"
	"third/gin"
	"time"
)

const (
	SAMPLE_RATE        = 1.0
	RPT_INTERVAL       = 60 * time.Second
	STATS_RPT_CHAN_BUF = 1e5
)

var (
	_statter      g2s.Statter
	_statsAggrMap map[string]*StatsRpt
	_statsRptCh   chan *StatsRpt
	_reg          *regexp.Regexp = regexp.MustCompile(`[.:/]`)
)

type StatsRpt struct {
	Bucket   string
	Count    int
	Dur      time.Duration
	CreateAt time.Time
}

func initStatsD(addr string) {
	if addr == "" {
		log.Printf("initStatsD, addr is empty")
		os.Exit(1)
	}

	statter, err := g2s.Dial("udp", addr)
	if err != nil {
		log.Println("initStatsD, %s error:%v", addr, err)
		os.Exit(1)
	}
	_statter = statter
	_statsAggrMap = make(map[string]*StatsRpt)
	_statsRptCh = make(chan *StatsRpt, STATS_RPT_CHAN_BUF)
	go consumeStats()
}

func Counter(bucket string, n ...int) {
	_statter.Counter(SAMPLE_RATE, bucket, n...)
}

func Timing(bucket string, d ...time.Duration) {
	_statter.Timing(SAMPLE_RATE, bucket, d...)
}

func Gauge(bucket string, v ...string) {
	_statter.Gauge(SAMPLE_RATE, bucket, v...)
}

func bucketName(sn, host, path, method string, httpCode int) string {
	host = _reg.ReplaceAllString(host, "_")
	path = _reg.ReplaceAllString(path, "_")
	return fmt.Sprintf("%s.%s.%d.%s_%s", sn, host, httpCode, path, method)

}

func goroutineStats(serviceName string) {
	hostName, _ := os.Hostname()
	hostName = _reg.ReplaceAllString(hostName, "_")
	bucket := fmt.Sprintf("%s.%s.goroutine", serviceName, hostName)
	for {
		n := runtime.NumGoroutine()
		Gauge(bucket, fmt.Sprintf("%d", n))
		time.Sleep(5 * time.Second)
	}
}

func gcStats(serviceName string) {
	hostName, _ := os.Hostname()
	hostName = _reg.ReplaceAllString(hostName, "_")
	bucket := fmt.Sprintf("%s.%s.gc", serviceName, hostName)
	gcs := &debug.GCStats{}
	for {
		debug.ReadGCStats(gcs)
		Gauge(bucket, fmt.Sprintf("%d", gcs.NumGC))
		time.Sleep(5 * time.Second)
	}
}

func consumeStats() {
	lastRptTime := time.Now()
	for {
		select {
		case s := <-_statsRptCh:
			rpt, ok := _statsAggrMap[s.Bucket]
			if !ok {
				_statsAggrMap[s.Bucket] = s
			} else {
				rpt.Count += s.Count
				rpt.Dur += s.Dur
			}
			if time.Now().Sub(lastRptTime) > RPT_INTERVAL {
				for _, rpt := range _statsAggrMap {
					Counter(rpt.Bucket, rpt.Count)
					Timing(rpt.Bucket, rpt.Dur)
				}
				lastRptTime = time.Now()
				_statsAggrMap = make(map[string]*StatsRpt)
				// log.Printf("flush all statsd")
			}
			// log.Printf("statsAggrMap:%+v", _statsAggrMap)
		}
	}
}

func GinStatter(statsdAddr, serviceName string) gin.HandlerFunc {
	initStatsD(statsdAddr)
	serviceName = _reg.ReplaceAllString(serviceName, "_")

	go goroutineStats(serviceName)
	go gcStats(serviceName)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		bucket := bucketName(serviceName, c.Request.Host, c.Request.URL.Path, c.Request.Method, c.Writer.Status())
		// log.Printf("statsd bucket:%s", bucket)
		duration := time.Now().Sub(start)
		_statsRptCh <- &StatsRpt{
			Bucket:   bucket,
			Count:    1,
			Dur:      duration,
			CreateAt: time.Now(),
		}
	}
}
