package subcmd

import (
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/unix1/blug/internal/site"
)

const (
	servePortTries = 5
	ephemeralMin   = 49152
	ephemeralMax   = 65535
)

func Serve(siteRoot string) error {
	publicDir := filepath.Join(siteRoot, site.PublicDir)
	if _, err := os.Stat(publicDir); err != nil {
		return fmt.Errorf("%s/: %w", site.PublicDir, err)
	}
	ln, err := listenHTTP()
	if err != nil {
		return err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	fmt.Printf("Serving %s at http://localhost:%d\n", publicDir, port)
	return http.Serve(ln, http.FileServer(http.Dir(publicDir)))
}

func listenHTTP() (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", site.DefaultServePort))
	if err == nil {
		return ln, nil
	}
	last := err
	for range servePortTries {
		port := ephemeralMin + rand.IntN(ephemeralMax-ephemeralMin+1)
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			return ln, nil
		}
		last = err
	}
	return nil, fmt.Errorf("listen: %w", last)
}
