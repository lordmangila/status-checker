package main

import (
	"github.com/lordmangila/status-checker/pkg/rest"
)

func main() {
	server := rest.NewServer()
	go server.Run()

	server.SetRoutes()

	server.ListenListenAndServe()
}
