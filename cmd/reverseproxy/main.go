// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	domain1 string
	proxy1recvr string

)

func init() {
	log.SetPrefix("quikReverseProxy: ")
}

type rProxy struct {
	addr     string
	port     net.Addr
	hostname string
}

func (p rProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Hostname() == domain1 {
		proxy := httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Host = proxy1recvr
				req.URL.Scheme = r.URL.Scheme
				req.Host = ""
			}
		}
	}
}
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// clienturl, err := url.Parse(r.RemoteAddr)
	// if err != nil {
	// 	log.Printf("error parsing client addr: %w\n", err)
	// 	panic(err)
	// }
	rp := httputil.NewSingleHostReverseProxy(r.URL)
	rp.ServeHTTP(w, r)

}

func ListenAndProxy(laddr string) (err error) {
	http.HandleFunc("/", handleRedirect)
	for {
		err = http.ListenAndServe(laddr, nil)
		if err != nil {
			break
		}

	}
	return err
}

func main() {
	log.Fatal(ListenAndProxy(os.Args[1]))
}
