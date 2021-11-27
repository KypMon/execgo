package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var levelMap = make(map[int]string)

func main() {
	// initialize log level
	levelMap[0] = "[SYSTEM] "
	levelMap[1] = "[INFO] "
	levelMap[2] = "[WARNING] "
	levelMap[3] = "[ERROR] "

	// HANDLER
	http.HandleFunc("/healthz", healthz)

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

func healthz(w http.ResponseWriter, r *http.Request) {
	// 4. return 200 for healthz
	var statusCode = http.StatusOK

	setHeader(w, r, statusCode)

	w.WriteHeader(statusCode)
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

	io.WriteString(w, "welcome to healthz!")
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
