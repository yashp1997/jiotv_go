package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/jchv/go-webview2"
)

// StartGUI starts the JioTV Go server in a background goroutine and opens a WebView2 window.
func StartGUI() {
	// Find a free port
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("Could not find a free port:", err)
	}
	port := conn.Addr().(*net.TCPAddr).Port
	conn.Close()

	host := "127.0.0.1"
	portStr := strconv.Itoa(port)
	url := fmt.Sprintf("http://%s:%s", host, portStr)

	// WaitGroup to ensure server is started before loading window (optional, but good practice)
	var wg sync.WaitGroup
	wg.Add(1)

	// Start Server in Goroutine
	go func() {
		// Initialize Config if not already done (Main does it, but purely for safety)
		// Assuming LoadConfig was called in main before this.

		log.Printf("Starting Server on %s via GUI mode...", url)
		app, err := JioTVServer(JioTVServerConfig{
			Host: host,
			Port: portStr,
			TLS:  false,
		})
		if err != nil {
			log.Fatal("Failed to create server:", err)
		}

		wg.Done()
		if err := app.Listen(fmt.Sprintf("%s:%s", host, portStr)); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	// Launch WebView
	debug := os.Getenv("DEBUG") == "true"
	
	w := webview2.New(debug)
	if w == nil {
		log.Fatal("Failed to load WebView2. Is the runtime installed?")
	}
	defer w.Destroy()

	w.SetTitle("JioTV Go")
	w.SetSize(1280, 720, webview2.HintNone)
	w.Navigate(url)
	w.Run()
}
