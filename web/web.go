package web

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Scheme string

const (
	SCHEME_HTTP        Scheme = "http"
	SCHEME_HTTPS       Scheme = "https"
	SCHEME_FCGI        Scheme = "fcgi"
	SCHEME_UNIX_SOCKET Scheme = "unix"
)

var (
	Protocol             Scheme = SCHEME_HTTP
	Domain               string = "localhost"
	HTTPAddr             string = "localhost"
	HTTPPort             string = "8081"
	AppURL               string = "http://localhost:8081"
	AppSubURL            string = ""
	CertFile, KeyFile    string
	UnixSocketPermission uint32
)

func Run(handler http.Handler, startPage string) {
	listenAddr := GetListenAddr()
	OpenStartPage(listenAddr, startPage)
	ListenAndServe(listenAddr, handler)
}

func GetListenAddr() string {
	var listenAddr string
	if Protocol == SCHEME_UNIX_SOCKET {
		listenAddr = fmt.Sprintf("%s", HTTPAddr)
	} else {
		listenAddr = fmt.Sprintf("%s:%s", HTTPAddr, HTTPPort)
	}
	log.Printf("Listen: %v://%s%s", Protocol, listenAddr, AppSubURL)
	return listenAddr
}

// Start listen and serve
func ListenAndServe(listenAddr string, r http.Handler) error {
	var err error
	switch Protocol {
	case SCHEME_HTTP:
		err = http.ListenAndServe(listenAddr, r)
	case SCHEME_HTTPS:
		server := &http.Server{Addr: listenAddr, TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS10,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // Required for HTTP/2 support.
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}, Handler: r}
		err = server.ListenAndServeTLS(CertFile, KeyFile)
	case SCHEME_FCGI:
		err = fcgi.Serve(nil, r)
	case SCHEME_UNIX_SOCKET:
		os.Remove(listenAddr)

		var listener *net.UnixListener
		listener, err = net.ListenUnix("unix", &net.UnixAddr{Name: listenAddr, Net: "unix"})
		if err != nil {
			break // Trigger error after switch
		}

		// FIXME: add proper implementation of signal capture on all protocols
		// execute this on SIGTERM or SIGINT: listener.Close()
		if err = os.Chmod(listenAddr, os.FileMode(UnixSocketPermission)); err != nil {
			log.Fatalf("Failed to set permission of unix socket: %v", err)
		}
		err = http.Serve(listener, r)
	default:
		log.Fatalf("Invalid protocol: %s", Protocol)
	}
	if err != nil {
		log.Fatalf("Fail to start server: %v", err)
	}
	return err
}

func OpenStartPage(listenAddr string, startPage string) {
	if len(startPage) == 0 {
		startPage = AppSubURL
	}
	if len(startPage) == 0 {
		return
	}

	if !strings.HasPrefix(startPage, "/") {
		startPage = "/" + startPage
	}
	err := OpenBrower(fmt.Sprintf("%v://%s%s", Protocol, listenAddr, startPage))
	if err != nil {
		log.Printf("open %s: %s", listenAddr, err)
	}
}

//open os default webbrower
func OpenBrower(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
