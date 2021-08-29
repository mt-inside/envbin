package main

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/mt-inside/go-usvc"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/mt-inside/envbin/internal/config"
	"github.com/mt-inside/envbin/internal/servers"
	"github.com/mt-inside/envbin/pkg/middleware"
)

var Serve = &cli.Command{
	Name:  "serve",
	Usage: "Serve data over the network",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "output-addr",
			Value: ":8080",
			Usage: "Listen address for Lorem Ipsum output",
		},
		&cli.StringFlag{
			Name:  "api-addr",
			Value: ":8081",
			Usage: "Listen address for API",
		},
	},

	Action: serve,
}

func serve(c *cli.Context) error {
	var err error

	log := c.App.Metadata["log"].(logr.Logger)

	stopCh := make(chan struct{})
	shutdown := false
	defer func() {
		if !shutdown {
			close(stopCh)
		}
	}()

	signalCh := usvc.InstallSignalHandlers(log)

	err = config.LoadConfig(log)
	if err != nil {
		panic(err)
	}

	// TODO: output does need to be on a different port cause it needs a different http server cause of low-level tcp-rst stuff
	muxApi := gin.Default()
	servers.GetProbes(log, muxApi.Group("/health"))
	servers.GetEnv(log, muxApi.Group("/api/v1/env"))
	servers.GetConfig(log, muxApi.Group("/api/v1/config"))

	listenAddrApi := c.String("api-addr")
	chApi := serveHttpSimple(log, listenAddrApi, muxApi, stopCh)

	// TODO: kick off ReadyTimer (as per -> ). 100% shouldn't be in config/

	muxOutput := gin.Default()
	servers.GetOutput(log, muxOutput.Group("/"))

	listenAddrOutput := c.String("output-addr")
	chOutput := serveHttpChaos(log, listenAddrOutput, muxOutput, stopCh)

	select {
	case err = <-chApi:
	case err = <-chOutput:
	case <-signalCh:
	}

	// TODO: keep the probes serving? Closing the http socket will be notice enough, but that's subject to a timeout, and we want /ready to fail instantly.
	// TODO: can we unmount the other routes?
	log.Info("Finished serving", "error", err)

	shutdown = true
	close(stopCh)

	shutdownDelay, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(viper.GetInt("ShutdownDelayMS")))
	defer cancel()

	log.Info("Waiting for shutdown")
	for _, ch := range []<-chan error{chOutput, chApi} {
		if err := <-ch; err != nil {
			log.Error(err, "Error during shutdown")
		}
	}

	log.Info("Delaying exit for the remainder of ShutdownDelay")
	<-shutdownDelay.Done()

	log.Info("Shutdown complete")

	return err
}

func serveHttpSimple(log logr.Logger, listenAddr string, router http.Handler, stopCh <-chan struct{}) <-chan error {
	localCh := make(chan error)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go serveHttpInner(log, srv, localCh, stopCh)

	return localCh
}

func serveHttpChaos(log logr.Logger, listenAddr string, router *gin.Engine, stopCh <-chan struct{}) <-chan error {
	localCh := make(chan error)

	// Order important; see comments
	router.Use(gin.WrapF(middleware.Delay)) // delay should apply to any return value
	router.Use(gin.WrapF(middleware.Error)) // error should come after a delay
	// Currently throttled at string-read time :( router.Use(gin.WrapH(middleware.Rate))  // don't apply rate limit to eg the "injected error" message
	// TODO: rps middleware, distinct from bw throttling above, eg https://github.com/s12i/gin-throttle

	deriveContext := func(ctx context.Context, c net.Conn) context.Context {
		return context.WithValue(
			context.WithValue(ctx,
				middleware.CtxKeyConn, c), // So the handler can close the connection at the TCP level. This is not the listen socket; it's called every time there's a new connection
			middleware.CtxKeyLog, log,
		)
	}
	connEvent := func(conn net.Conn, event http.ConnState) {
		log.V(1).Info("Connection event", "event", event, "remote", conn.RemoteAddr())

		if event == http.StateNew &&
			viper.GetString("ErrorType") == "tcp-rst" &&
			viper.GetFloat64("ErrorRate") > rand.Float64() {

			log.V(1).Info("Closing TCP stream with RST")

			// TODO: wait for delay period

			// close stream with bytes still to be read in the buffer, hence cause data loss, resulting in a RST
			conn.Close()
		}
	}
	srv := &http.Server{
		Addr:        listenAddr,
		Handler:     router,
		ConnContext: deriveContext,
		ConnState:   connEvent,
	}

	go serveHttpInner(log, srv, localCh, stopCh)

	return localCh
}

func serveHttpInner(log logr.Logger, srv *http.Server, localCh chan<- error, stopCh <-chan struct{}) {
	defer close(localCh)

	log.Info("Listening", "addr", srv.Addr)
	serverCh := usvc.ChannelWrapper(func() error { return srv.ListenAndServe() })

	select {
	case err := <-serverCh:
		localCh <- err
	case <-stopCh:
		// TODO: ideally reach out to all in-flight requests (like long-lived ones) and cancel them

		log.Info("Attempting to gracefully shut down http server")
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(viper.GetInt64("ShutdownDelayMS"))) // We use this value as the maximum and minimum time to wait
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error(err, "Failed to gracefully shut down http server")
			localCh <- err
		}
		log.Info("Http server shut down")
	}
}
