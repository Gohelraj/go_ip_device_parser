package go_ip_device_parser

import (
	"fmt"
	"net/http"

	ua "github.com/mileusna/useragent"
)

type UserAgentAndIPDetails struct {
	IP    string
	Agent struct {
		Browser struct {
			Name    string
			Version string
		}
		Device struct {
			Name string
		}
		Os struct {
			Name    string
			Version string
		}
	}
	IsMobile  bool
	IsTablet  bool
	IsDesktop bool
	IsBot     bool
}

// Generates device information including browser, device, os
func ParseUserAgentAndClientIP(r *http.Request) UserAgentAndIPDetails {
	var userAgent = UserAgentAndIPDetails{}

	// parse user agent string and return struct with filled details
	client := ua.Parse(r.UserAgent())

	if client.Device == "" && client.OS != "" {
		userAgent.Agent.Device.Name = fmt.Sprintf("%s System", client.OS)
	} else {
		userAgent.Agent.Device.Name = client.Device
	}

	// set struct values of agent in "UserAgentAndIPDetails" struct
	userAgent.Agent.Browser.Name = client.Name
	userAgent.Agent.Browser.Version = client.Version
	userAgent.Agent.Os.Name = client.OS
	userAgent.Agent.Os.Version = client.OSVersion
	userAgent.IsMobile = client.Mobile
	userAgent.IsTablet = client.Tablet
	userAgent.IsDesktop = client.Desktop
	userAgent.IsBot = client.Bot

	// get IP address from request
	ipAddress := GetClientIPAddress(r)

	// set IP address in "UserAgentAndIPDetails" struct
	userAgent.IP = ipAddress

	return userAgent
}