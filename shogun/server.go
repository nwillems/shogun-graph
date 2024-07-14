package shogun

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO: Register shutdown hooks
type Server struct {
	Server           *http.Server
	ShutdownTimeout  time.Duration
	shutdownFinished chan struct{}
	Log              *log.Logger
}

func (s *Server) ListenAndServe() error {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{}, 1)
	}

	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("unexpected error while listening. %w", err)
	}

	s.Log.Println("Waiting for shutdown finish")
	<-s.shutdownFinished
	s.Log.Println("Shutdown finished")

	return nil
}

func (s *Server) WaitForExit(ctx context.Context) {
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	<-waiter

	ctxShutdown, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
	defer cancel()

	err := s.Server.Shutdown(ctxShutdown)
	if err != nil {
		s.Log.Println("Shutting down: %w", err)
	} else {
		s.Log.Println("Shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
