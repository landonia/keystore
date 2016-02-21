// Landon Wainwright.

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/landonia/keystore"
	"github.com/landonia/keystore/transport"
)

// main will bootstrap the in-memory key/value store
func main() {

	// Define flags
	var httpAddr, tcpAddr, udpAddr, dataPath string
	flag.StringVar(&httpAddr, "httpAddr", ":8080", "the host:port to bind the HTTP server")
	flag.StringVar(&tcpAddr, "tcpAddr", ":8081", "the host:port to bind the TCP server")
	flag.StringVar(&udpAddr, "udpAddr", ":8082", "the host:port to bind the UDP server")
	flag.StringVar(&dataPath, "dataPath", "", "the path to the file for saving the key store")
	flag.Parse()

	// The program will run until it receives the correct signal
	done := GetSignalChannel()

	// Create a key store in disk
	keystore := keystore.NewService(dataPath)

	// Bind the network protocols that are required
	transport.StartHTTPServer(httpAddr, keystore.RequestChannel)
	transport.StartTCPServer(tcpAddr, keystore.RequestChannel)
	transport.StartUDPServer(udpAddr, keystore.RequestChannel)

	// Start
	keystore.Start()

	// Just wait to exit
	<-done
	<-keystore.Stop()
}

// GetSignalChannel will wait for an exit signal and send a complete flag
// on the returned boolean channel to indicate when to shutdown
func GetSignalChannel() <-chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	return done
}
