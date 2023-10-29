package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/servernoj/gobook/ch07/shortener"
	"github.com/servernoj/gobook/ch07/sqlx"
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

	db, err := sqlx.Dial()
	if err != nil {
		log.Fatalf("unable to dial DB: %s", err)
	}

	service := shortener.NewService(db)

	appServer := shortener.Server{}
	appServer.Init(service)

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
