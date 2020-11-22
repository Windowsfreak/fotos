package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"fotos/domain"
	"fotos/endpoint/pictures"
	"fotos/repository"
)

// Run starts the HTTP server
func Run() {
	log.SetFormatter(&log.TextFormatter{})

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&domain.Config)
	if err != nil {
		log.Fatal(err)
	}

	if len(domain.Config.PreSharedKey) < 1 {
		log.Fatal("pre-shared key missing in config")
	}

	handler := setUpServer()
	srv := &http.Server{Addr: domain.Config.ServerAddr, Handler: handler}
	go func() {
		log.WithField("server_port", srv.Addr).Info("Starting server")

		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Info("Server shut down. Waiting for connections to drain.")
			} else {
				log.WithError(err).
					WithField("server_port", srv.Addr).
					Fatal("failed to start server")
			}
		}
	}()

	// Wait for an interrupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)    // interrupt signal sent from terminal
	signal.Notify(sigint, syscall.SIGTERM) // sigterm signal sent from system
	<-sigint

	log.Info("Shutting down server")

	attemptGracefulShutdown(srv)
}

func setUpServer() http.Handler {
	mux := http.NewServeMux()

	picturesLogger := log.WithField("endpoint", "fotos")
	picturesRepository, err := repository.NewRepository()
	if err != nil {
		picturesLogger.Fatal(err)
	}
	picturesService := pictures.NewService(picturesLogger, picturesRepository)
	mux.Handle("/pictures/add", pictures.MakeAddPictureHandler(picturesService, picturesLogger))
	mux.Handle("/pictures/del", pictures.MakeDeletePictureHandler(picturesService, picturesLogger))
	mux.Handle("/pictures/random", pictures.MakeGetRandomPictureHandler(picturesService, picturesLogger))

	return mux
}

func attemptGracefulShutdown(srv *http.Server) {
	if err := shutdownServer(srv, 25*time.Second); err != nil {
		log.WithError(err).Error("failed to shutdown server")
	}
}

func shutdownServer(srv *http.Server, maximumTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), maximumTime)
	defer cancel()
	return srv.Shutdown(ctx)
}
