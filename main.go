package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	"%{time} %{level:.5s} %{message}",
)

func init() {
	stdout := logging.NewLogBackend(os.Stdout, "app ", 0)
	stdoutFormatter := logging.NewBackendFormatter(stdout, format)
	stdoutLeveled := logging.AddModuleLevel(stdout)
	stdoutLeveled.SetLevel(logging.INFO, "")
	logging.SetBackend(stdoutFormatter)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}

	name = strings.ToLower(name)

	output := RandStringBytes(10)

	switch name {
	case "debug":
		logger.Debug(output)
	case "info":
		logger.Info(output)
	case "warn":
		logger.Warning(output)
	case "error":
		logger.Error(output)
	default:
		logger.Info(output)

	}

	//logger.Info("Received request for %s\n", name)
	//w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	address := "0.0.0.0:8000"

	srv := &http.Server{
		Handler: r,
		Addr:    address,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 10 * time.Second,
	}

	// Start Server
	//go func() {
	//	logger.Info("Starting Server", address)
	//	if err := srv.ListenAndServe(); err != nil {
	//		logger.Fatal(err)
	//	}
	//}()
	logger.Info("Starting Server", address)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}

	// Graceful Shutdown
	//waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	logger.Error("Shutting down")
	os.Exit(0)
}
