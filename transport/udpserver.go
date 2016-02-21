// Landon Wainwright.

package transport

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/landonia/keystore"
)

// UDPServer holds the UDP client connection
type UDPServer struct {
	addr      string                   // the address to bind to
	requests  chan<- *keystore.Request // The request event channel to send the requests
	udpconn   *net.UDPConn             // The udp connection
	quit      chan bool                // The channel to wait on to finish the connection
	connected bool                     // Whether the server is currently connected
}

// StartUDPServer will start a new UDP service allowing requests
// to be made to the key store service
func StartUDPServer(addr string, requests chan<- *keystore.Request) {

	// Start the udp server
	go func() {

		// Create the server and start it up
		log.Printf("Starting UDP server using address: %s", addr)
		server := newUDPServer(addr, requests)
		server.connect()
	}()
}

// newUDPServer will create a new UDP server using the address
func newUDPServer(addr string, requests chan<- *keystore.Request) *UDPServer {
	return &UDPServer{addr: addr, requests: requests, quit: make(chan bool)}
}

// connect will enable the event listener for incoming data
func (server *UDPServer) connect() {
	if server.connected {
		log.Printf("The UDP server is already connected using address: %s", server.addr)
		return
	}

	// Now create a connection
	udpAddr, err := net.ResolveUDPAddr("udp", server.addr)
	if err == nil {
		server.udpconn, err = net.ListenUDP("udp", udpAddr)
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("UDP server now connected to address: %s", server.addr)
	server.connected = true

	go func() {
		buf := make([]byte, 1024)
		for {

			// Collect the bytes from the socket
			n, client, err := server.udpconn.ReadFromUDP(buf)
			if err != nil {
				log.Println("Error whilst reading UDP packet: ", err)
			} else {

				// Attempt to read the data into a request
				request := &keystore.Request{}
				clientaddr := client.String()
				if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(request); err != nil {
					log.Printf("Error whilst reading UDP packet: %s", err)

					// Send an error response back to the client
					go func() {

						// The channel can not be sent so will be created
						// Now we need to send the response back to the client
						log.Printf("Sending UDP error response to client [%s]", clientaddr)

						// Handle the response
						writeResponse(server.udpconn, client, &keystore.Response{Error: err.Error()})
					}()
				} else {
					log.Printf("Received UDP request from client: [%s]", clientaddr)

					// Wait for the response and send it back to the client
					go func() {

						// The channel can not be sent so will be created
						request.ResponseChannel = make(chan *keystore.Response)

						// Now send the request on the request channel
						server.requests <- request
						response := <-request.ResponseChannel

						// Now we need to send the response back to the client
						log.Printf("Received response... Sending UDP response to client [%s]", clientaddr)

						// Handle the response
						writeResponse(server.udpconn, client, response)
					}()
				}
			}
		}
	}()
}

// writeResponse will write the response object back to the udp client
func writeResponse(udpconn *net.UDPConn, client *net.UDPAddr, response *keystore.Response) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(response); err != nil {
		log.Printf("Error writing UDP response to buffer: %s", err)
	}

	// Write the response to the client
	_, err := udpconn.WriteToUDP(buf.Bytes(), client)
	if err != nil {
		log.Printf("Error writing UDP response to client: %s", err)
	}
}
