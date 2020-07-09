package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/gargath/metrics-example/pkg/server"
)

func main() {
	log.Printf("Example API %s\n", version())

	viper.SetEnvPrefix("METRICSEXAMPLE")
	viper.AutomaticEnv()

	flag.String("listenAddr", "0.0.0.0:8080", "address to listen on")
	flag.Bool("help", false, "print this help and exit")

	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		panic(err)
	}

	if viper.GetBool("help") {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprint(os.Stderr, flag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	s := &server.Server{
		GracefulShutdownPeriod: time.Second * 60,
		Addr:                   viper.GetString("listenAddr"),
	}
	if err := runServer(s); err != nil {
		log.Printf("startup failed: %v", err)
	}
}

func runServer(s *server.Server) error {
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")
		if err := s.Shutdown(); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	err := s.Run()
	if err != nil {
		return err
	}

	<-done
	log.Println("Server stopped")
	return nil
}
