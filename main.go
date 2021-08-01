package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	svr := &http.Server{Addr: ":9090"}
	// http server
	g.Go(func() error {
		go func() {
			<-ctx.Done()
			svr.Shutdown(ctx)
		}()
		return svr.ListenAndServe()
	})

	// signal
	g.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
		sig := make(chan os.Signal, len(exitSignals))
		signal.Notify(sig, exitSignals...)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-sig:
				return nil
			}
		}
	})

	// inject error
	g.Go(func() error {
		time.Sleep(time.Second)
		return errors.New("inject error")
	})

	err := g.Wait() // first error return
	fmt.Println(err)
}
