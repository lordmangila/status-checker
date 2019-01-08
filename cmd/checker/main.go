package main

import (
	"bitbucket.org/lordmangila/status-checker/pkg/rest"
)

func main() {
	server := rest.NewServer()
	go server.Run()

	server.SetRoutes()

	server.ListenListenAndServe()
}
