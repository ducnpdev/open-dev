package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	Reuse()
	//
	NonReuse()
}

// NonReuse, not reuse http
func NonReuse() {
	// client trace to log whether the request's underlying tcp connection was re-used
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("conn was reused: %t", info.Reused)
			log.Printf("%#v", info)
		},
		DNSDone: func(di httptrace.DNSDoneInfo) {
			log.Printf("DNSDone: %#v", di)
		},
		DNSStart: func(di httptrace.DNSStartInfo) {
			log.Printf("DNSStart: %#v", di)
		},
		ConnectStart: func(network, addr string) {
		},
		Wait100Continue: func() {},
	}
	traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)

	// 1st request
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// 2nd request
	req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}

// Reuse, reuse http client
func Reuse() {
	var (
		err error
	)
	// client trace to log whether the request's underlying tcp connection was re-used
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("conn was reused: %t", info.Reused)
			log.Printf("%#v", info)
		},
	}
	traceCtx := httptrace.WithClientTrace(context.Background(), clientTrace)

	// 1st request
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	// 2nd request
	req, err = http.NewRequestWithContext(traceCtx, http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
