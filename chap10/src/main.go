package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var levelMap = make(map[int]string)

func main() {
	// initialize log level
	levelMap[0] = "[SYSTEM] "
	levelMap[1] = "[INFO] "
	levelMap[2] = "[WARNING] "
	levelMap[3] = "[ERROR] "

	// register metrics
	Register()

	// HANDLER
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())

	// provide an default env variable
	_, isLogLevelPresent := os.LookupEnv("LOG_LEVEL")
	if !isLogLevelPresent {
		os.Setenv("LOG_LEVEL", "0")
	}

	_, isPortPresent := os.LookupEnv("PORT")
	if !isPortPresent {
		os.Setenv("PORT", "8080")
	}

	_, isVersionPresent := os.LookupEnv("VERSION")
	if !isVersionPresent {
		os.Setenv("VERSION", "1.4.0")
	}

	//fmt.Println("Listening " + os.Getenv("PORT") + ": ", 2)

	logger("Listening "+os.Getenv("PORT")+": ", 1)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func logger(msg string, logLevel int) {

	lvl, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	// if this message log level is higher/more verbose than the preset level, then output
	if logLevel >= lvl {
		fmt.Println(levelMap[logLevel] + msg)
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	// #1: add delay here
	delay := randInt(0, 2000)
	fmt.Printf("Sleeping %d milliseconds...\n", delay)
	time.Sleep(time.Millisecond * time.Duration(delay))

	var statusCode = http.StatusOK
	setHeader(w, r, statusCode)
	w.WriteHeader(statusCode)

	io.WriteString(w, "welcome to home page!")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	// 4. return 200 for healthz
	var statusCode = http.StatusOK

	setHeader(w, r, statusCode)

	w.WriteHeader(statusCode)

	io.WriteString(w, "welcome to healthz!")
}

func setHeader(w http.ResponseWriter, r *http.Request, statusCode int) {
	// 1. set req header to resp header
	for name, headers := range r.Header {
		for _, h := range headers {
			w.Header().Set(name, h)
		}
	}

	// 2. set version as the header
	w.Header().Set("VERSION", os.Getenv("VERSION"))

	ip, err := getIP(r)
	//fmt.Println(ip)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No valid ip"))
	}

	// 3 record ip and status code
	logger("the IP of this request is: "+ip+"The Status of the request is: "+strconv.Itoa(statusCode), 1)

	//io.WriteString(w, "welcome to healthz!")
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

const (
	MetricsNamespace = "httpserver"
)

// NewExecutionTimer provides a timer for Updater's RunOnce execution
func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

var (
	functionLatency = CreateExecutionTimeMetric(MetricsNamespace,
		"Time spent.")
)

// NewExecutionTimer provides a timer for admission latency; call ObserveXXX() on it to measure
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

// ObserveTotal measures the execution time from the creation of the ExecutionTimer
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

// CreateExecutionTimeMetric prepares a new histogram labeled with execution step
func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

// ExecutionTimer measures execution time of a computation, split into major steps
// usual usage pattern is: timer := NewExecutionTimer(...) ; compute ; timer.ObserveStep() ; ... ; timer.ObserveTotal()
type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}
