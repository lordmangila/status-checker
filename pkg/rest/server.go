package rest

import (
	"flag"
	"log"
	"net/http"
	"time"

	"bitbucket.org/lordmangila/status-checker/pkg/checker"

	"github.com/gorilla/websocket"
)

// Time to re-check url status.
const checkTime = 5 * time.Second

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Server contains information about the server.
type Server struct {
	// Registered client list.
	clients map[*Client]bool
}

// NewServer initializes an new websocket server.
func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
	}
}

// Run ...
func (s *Server) Run() {
	for {
		for client := range s.clients {
			client.CheckSites()
			time.Sleep(checkTime)
		}
	}
}

// SetRoutes sets the routes.
func (s *Server) SetRoutes() {
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/check", s.ServeWS)
}

// ListenListenAndServe initiates the handlers.
func (s *Server) ListenListenAndServe() {
	flag.Parse()

	http.ListenAndServe(*addr, nil)
}

// ServeHome serves the home page file handler.
func ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/main.gohtml")
}

// ServeWS serves the websocket handler.
func (s *Server) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		log.Println(err)
		conn.Close()
	}

	client := &Client{
		server: s,
		sites:  make(map[*checker.Site]bool),
		conn:   conn,
		send:   make(chan []byte, 256),
	}
	s.clients[client] = true

	go client.Broadcast()
	go client.Listen()
}
