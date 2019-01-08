package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/gorilla/websocket"

	"bitbucket.org/lordmangila/status-checker/pkg/rest"
)

const (
	// UpgradeWebsocket defines the Upgrade websocket header value
	UpgradeWebsocket = "websocket"

	// ConnectionUpgrade defines the Connection Upgrade header value
	ConnectionUpgrade = "Upgrade"
)

var rs *rest.Server

func init() {
	// filename points to ./status-checker/pkg/rest/server_test.go.
	_, filename, _, _ := runtime.Caller(0)

	// Move up to status-checker root directory.
	dir := path.Join(path.Dir(filename), "../..")

	// Change the current working directory to project root for ServeHome.
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	// Start the rest server.
	rs = rest.NewServer()
	go rs.Run()

	rs.SetRoutes()

	go func() {
		rs.ListenListenAndServe()
	}()
}

func ExampleServer() {
	// Create a new test Server.
	server := httptest.NewServer(http.HandlerFunc(rs.ServeWS))
	defer server.Close()

	// Connect to the websocket address.
	wsAddr := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if nil != err {
		conn.Close()
	}

	// Send an invalid url to the server.
	conn.WriteMessage(1, []byte("invalidurl"))

	// Read failed response from the websocket.
	_, message, _ := conn.ReadMessage()
	fmt.Println(string(message))

	// Register https://www.google.com/ to the client's sites.
	conn.WriteMessage(1, []byte("https://www.google.com/"))

	// Read response from the websocket for google.
	_, message, _ = conn.ReadMessage()
	fmt.Println(string(message))

	// Register https://www.github.com/ to the client's sites.
	conn.WriteMessage(1, []byte("https://www.github.com/"))

	// This will be done twice as there are already 2 sites registered to the client.
	// Read response from the websocket.
	// 1st reponse will show google status.
	_, message, _ = conn.ReadMessage()
	fmt.Println(string(message))

	// Read response from the websocket.
	// 2nd reponse will show github status.
	_, message, _ = conn.ReadMessage()
	fmt.Println(string(message))

	// output:
	// {"URL":"invalidurl","StatusCode":0,"Active":false,"Valid":false,"Error":"Invalid URI: invalidurl"}
	// {"URL":"https://www.google.com/","StatusCode":200,"Active":true,"Valid":true,"Error":""}
	// {"URL":"https://www.google.com/","StatusCode":200,"Active":true,"Valid":true,"Error":""}
	// {"URL":"https://www.github.com/","StatusCode":200,"Active":true,"Valid":true,"Error":""}
}

func TestServeHome(t *testing.T) {
	// Create a new Request to be passed to the handler.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test ResponseRecorder to record the response and pass
	// to the handler.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rest.ServeHome)

	// Call ServeHTTP and pass Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check ResponseRecorder status code returns 200.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Invalid status code: got %v want %v", status, http.StatusOK)
	}
}

func TestServeWS(t *testing.T) {
	// Create a new test Server.
	server := httptest.NewServer(http.HandlerFunc(rs.ServeWS))
	defer server.Close()

	// Connect to the websocket address.
	wsAddr := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, resp, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if nil != err {
		t.Logf("Error: %v", err)
		conn.Close()
	}

	// Check ResponseRecorder response header Upgrade conforms with websocket headers.
	if ug := resp.Header["Upgrade"][0]; UpgradeWebsocket != ug {
		t.Errorf("Invalid Upgrade header: got %v want %v", ug, UpgradeWebsocket)
	}

	// Check ResponseRecorder response header Connection conforms with websocket headers.
	if c := resp.Header["Connection"][0]; ConnectionUpgrade != c {
		t.Errorf("Invalid Upgrade header: got %v want %v", c, ConnectionUpgrade)
	}
}
