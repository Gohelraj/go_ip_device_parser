package go_ip_device_parser

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
)

type UserAgent struct {
	user_agent.UserAgent
	IpAddress string `json:"ipAddress"`
}

func ParseUserAgent(r *http.Request) *UserAgent {
	var userAgent = UserAgent{}
	ua := user_agent.New(r.UserAgent())
	userAgent.UserAgent = *ua

	ipAddress, _ := getIP(r)
	userAgent.IpAddress = ipAddress

	return &userAgent
}

func getIP(r *http.Request) (string, error) {
	var ipAddress string
	var netIP net.IP

	if r.Header != nil {
		ipAddress = r.Header.Get("x-client-ip")
		netIP = net.ParseIP(ipAddress)
		if netIP != nil {
			return ipAddress, nil
		}
	}

	// Get IP from the X-REAL-IP header
	ipAddress = r.Header.Get("cf-connecting-ip")
	netIP = net.ParseIP(ipAddress)
	if netIP != nil {
		return ipAddress, nil
	}

	// Get IP from the X-REAL-IP header
	ipAddress = r.Header.Get("true-client-ip")
	netIP = net.ParseIP(ipAddress)
	if netIP != nil {
		return ipAddress, nil
	}

	// Get IP from the X-REAL-IP header
	ipAddress = r.Header.Get("x-real-ip")
	netIP = net.ParseIP(ipAddress)
	if netIP != nil {
		return ipAddress, nil
	}

	// Get IP from the X-REAL-IP header
	ipAddress = r.Header.Get("x-forwarded")
	netIP = net.ParseIP(ipAddress)
	if netIP != nil {
		return ipAddress, nil
	}

	ipAddress = r.Header.Get("forwarded-for")
	netIP = net.ParseIP(ipAddress)
	if netIP != nil {
		return ipAddress, nil
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("x-forwarded-for")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	// Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}
