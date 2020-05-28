# go_ip_device_parser

Go package to Parse `IP Adress` and `Client Device Information` from the request.

## Installation

```
go get -u github.com/Gohelraj/go_ip_device_parser
```

## API

```js
package main

import (
	"fmt"
	"log"
	"net/http"
	
	ip_parser "github.com/Gohelraj/go_ip_device_parser"
)

func main() {
	http.HandleFunc("/", HelloServer)
	port := ":3004"
	fmt.Printf("Server listening on port: %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	clientDetails := ip_parser.ParseUserAgentAndClientIP(r)
    	fmt.Fprintln(w,"Browser Name: ", clientDetails.Agent.Browser.Name)
    	fmt.Fprintln(w,"Browser Version: ", clientDetails.Agent.Browser.Version)
    	fmt.Fprintln(w,"Device Name: ", clientDetails.Agent.Device.Name)
    	fmt.Fprintln(w,"OS Name: ", clientDetails.Agent.Os.Name)
    	fmt.Fprintln(w,"OS Version: ", clientDetails.Agent.Os.Version)
    	fmt.Fprintln(w,"Client IP: ", clientDetails.IP)
    	fmt.Fprintln(w,"IsBot: ", clientDetails.IsBot)
    	fmt.Fprintln(w,"IsDesktop: ", clientDetails.IsDesktop)
    	fmt.Fprintln(w,"IsMobile: ", clientDetails.IsMobile)
    	fmt.Fprintln(w,"IsTablet: ", clientDetails.IsTablet)
}

```

## How It Works

It looks for specific headers in the request and falls back to some defaults if they do not exist.

The following is the order we use to determine the user ip from the request.

1. `X-Client-IP`  
2. `X-Forwarded-For` (Header may return multiple IP addresses in the format: "client IP, proxy 1 IP, proxy 2 IP", so we take the the first one.)
3. `CF-Connecting-IP` (Cloudflare)
4. `True-Client-Ip` (Akamai and Cloudflare)
5. `X-Real-IP` (Nginx proxy/FastCGI)
6. `X-Cluster-Client-IP` (Rackspace LB, Riverbed Stingray)
7. `X-Forwarded`, `Forwarded-For` and `Forwarded` (Variations of #2)
8. `req.remoteAddress`

If an IP address cannot be found, it will return `""`.

## References
http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/x-forwarded-headers.html \
http://stackoverflow.com/a/11285650 \
http://www.squid-cache.org/Doc/config/forwarded_for/ \
https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers- \
http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
