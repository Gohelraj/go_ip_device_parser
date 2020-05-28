package go_ip_device_parser

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// Get client's IP address
func getIPAddress(r *http.Request) (string, error) {
	var ipAddress string
	var netIP net.IP

	if r.Header != nil {
		// Standard headers used by Amazon EC2, Heroku, and others.
		ipAddress = r.Header.Get("x-client-ip")
		netIP = net.ParseIP(ipAddress)
		if netIP != nil {
			return ipAddress, nil
		}

		// Cloudflare.
		// @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
		// CF-Connecting-IP - applied to every request to the origin.
		ipAddress = r.Header.Get("cf-connecting-ip")
		netIP = net.ParseIP(ipAddress)
		if netIP != nil {
			return ipAddress, nil
		}

		// Akamai and Cloudflare: True-Client-IP.
		ipAddress = r.Header.Get("true-client-ip")
		netIP = net.ParseIP(ipAddress)
		if netIP != nil {
			return ipAddress, nil
		}

		// Default nginx proxy/fcgi; alternative to x-forwarded-for, used by some proxies.
		ipAddress = r.Header.Get("x-real-ip")
		netIP = net.ParseIP(ipAddress)
		if netIP != nil {
			return ipAddress, nil
		}

		// (Rackspace LB and Riverbed's Stingray)
		// http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
		// https://splash.riverbed.com/docs/DOC-1926
		ipAddress = r.Header.Get("x-cluster-client-ip")
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

		// x-forwarded-for may return multiple IP addresses in the format:
		// "client IP, proxy 1 IP, proxy 2 IP"
		// Therefore, the right-most IP address is the IP address of the most recent proxy
		// and the left-most IP address is the IP address of the originating client.
		// source: http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/x-forwarded-headers.html
		// Azure Web App's also adds a port for some reason, so we'll only use the first part (the IP)
		// Load-balancers (AWS ELB) or proxies.
		ips := r.Header.Get("x-forwarded-for")
		splitIps := strings.Split(ips, ",")
		for _, ip := range splitIps {
			netIP = net.ParseIP(ip)
			if netIP != nil {
				return ip, nil
			}
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
