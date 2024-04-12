package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Downstream struct {
	URL string `yaml:"url"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Conf struct {
	Server     Server     `yaml:"server"`
	Downstream Downstream `yaml:"downstream"`
	Whitelist  []string   `yaml:"whitelist"`
}

var conf = Conf{}

func transformToCIDR(ipStr string) string {
	if ip := net.ParseIP(ipStr); ip != nil {
		if ip.To4() != nil {
			// This is IPv4
			return ip.String() + "/32"
		} else {
			// This is IPv6
			return ip.String() + "/128"
		}
	}
	return ipStr
}

func ipInWhitelist(ipStr string, whitelist []string) bool {
	ip := net.ParseIP(ipStr)
	for _, cidr := range whitelist {
		_, ipNet, err := net.ParseCIDR(transformToCIDR(cidr))
		if err != nil {
			return false
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Enter handler\n")
	ipStr := r.Header.Get("X-FORWARDED-FOR")
	if ipStr == "" {
		ipStr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		ipStr = strings.TrimSpace(strings.Split(ipStr, ",")[0])
	}
	log.Printf("ip: %s\n", ipStr)

	if ipInWhitelist(ipStr, conf.Whitelist) {
		serveReverseProxy(conf.Downstream.URL, w, r)
		log.Printf("ip in whitelist, exit handler, ipstr: %s\n", ipStr)
		return
	}

	http.NotFound(w, r)
	log.Printf("Exit handler\n")
}

func serveReverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}

func main() {
	log.SetOutput(os.Stdout)
	// Load configuration
	source, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Panic(err)
	}

	err = yaml.Unmarshal(source, &conf)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("start http server , listen on %s\n", conf.Server.Addr)
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(conf.Server.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
