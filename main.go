package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/crypto/acme/autocert"
)

var redirectPtr string

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix("http", redirectPtr) {
		redirectPtr = "https://" + redirectPtr
	}
	http.Redirect(w, r, redirectPtr, http.StatusFound)
}

func main() {
	var fakeDomains = flag.String("fake-domain", "", "Comma delimited list of domains to host")
	flag.StringVar(&redirectPtr, "redirect", "", "Domain to redirect to")
	var local = flag.Bool("local", false, "Local server for testing, do not get cert")
	flag.Parse()

	hosts := strings.Split(*fakeDomains, ",")

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...), //Your domain here
		Cache:      autocert.DirCache("."),           // Store certs where we execute from
	}

	http.HandleFunc("/", catchAllHandler)

	if *local {
		fmt.Println("AutoCert: False, TLS: False")
		err := http.ListenAndServe(":8443", nil)
		if err != nil {
			fmt.Printf("main(): %s\n", err)
		}
	} else {
		server := &http.Server{
			Addr: ":8443",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		err := server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt
		if err != nil {
			fmt.Printf("main(): %s\n", err)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown signal received, exiting fakeit")
}
