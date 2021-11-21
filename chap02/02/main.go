package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// provide an default env variable
	//os.Setenv("VERSION", "1.0.0")

	http.HandleFunc("/healthz", healthz)

	fmt.Println("Listening 8080:")

	http.ListenAndServe(":8080", nil)
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
	fmt.Println(ip)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No valid ip"))
	}

	// 3 record ip and status code
	fmt.Println("the IP of this request is: ", ip, "The Status of the request is: ", statusCode)
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
