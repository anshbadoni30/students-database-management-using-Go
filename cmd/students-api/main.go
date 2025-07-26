package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anshbadoni30/students-api/internal/config"
)

func homepagehandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Welcome to home page"))
}

func main() {
	//setup config
	cfg:= config.MustLoad()

	//setup database

	//setup routes
	router:= http.NewServeMux()
	router.HandleFunc("/",homepagehandler)

	//setup server
	server:=http.Server{
		Addr: cfg.HttpServer.Address,
		Handler: router,
	}
	slog.Info("server started", slog.String("on Address: ", cfg.HttpServer.Address))
	
	done:= make(chan os.Signal,1)
	signal.Notify(done, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){
		err:= server.ListenAndServe()
		if err!=nil{
			log.Fatal("failed to start the server")
		}
	} ()

	<-done

	slog.Info("shutting down the server")
	
	ctx, cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //Even if the timeout is hit, calling cancel() explicitly ensures that internal timers and goroutines are cleaned up properly.
	err:=server.Shutdown(ctx)
	if err!=nil{
		slog.Error("failed to shutdown server", slog.String("error",err.Error()))
	}
	slog.Info("server shutdown successfully")
	/*
	First, the server starts and begins listening for requests.
When the user presses Ctrl+C, the done channel receives the termination signal.
The <-done line unblocks, and context.WithTimeout() is called to create a ctx with a 5-second timeout.
This ctx is passed into server.Shutdown(ctx), giving the server up to 5 seconds to finish any in-flight requests.
After that, the server shuts down gracefully — either when requests complete or the timeout expires. */

}
