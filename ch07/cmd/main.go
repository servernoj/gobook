package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/servernoj/gobook/ch07/shortener"
)

type flags struct {
	addr string
}

func main() {
	f := flags{
		addr: ":8080",
	}
	flag.StringVar(&f.addr, "a", f.addr, "network address to listen")
	flag.Parse()

	logger := log.New(os.Stderr, "shortener: ", log.LstdFlags|log.Lmsgprefix)
	logger.Printf("starting the server on %s\n", f.addr)

	appServer := shortener.Server{}
	appServer.Init()

	httpServer := http.Server{
		Addr:     f.addr,
		Handler:  appServer,
		ErrorLog: logger,
	}
	defer httpServer.Close()
	if os.Getenv("USE_LOGGER") == "1" {
		httpServer.Handler = shortener.Logger(httpServer.Handler)
	}

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, "server closed unexpectedly:", err)
	}

}
