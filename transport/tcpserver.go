// Landon Wainwright.

package transport

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/landonia/keystore"
)

// TCPServer holds the TCP client connection
type TCPServer struct {
	addr      string                   // the address to bind to
	requests  chan<- *keystore.Request // The request event channel to send the requests
	listener  net.Listener             // The tcp connection
	quit      chan bool                // The channel to wait on to finish the connection
	connected bool                     // Whether the server is currently connected
}

// TCPClientHandler holds the TCP client connection
type TCPClientHandler struct {
	requestChannel chan<- *keystore.Request // The request channel
	conn           net.Conn                 // The tcp connection
	quit           chan bool                // The channel to wait on to finish the connection
	encoder        *gob.Encoder             // The encoder for this connection
	decoder        *gob.Decoder             // The decoder for this connection
}

// StartTCPServer will start a new TCP server allowing requests
// to be made to the key store service
func StartTCPServer(addr string, requests chan<- *keystore.Request) {

	// Start the server
	go func() {

		// Create the server and start it up
		log.Printf("Starting TCP server using address: %s", addr)
		server := newTCPServer(addr, requests)
		server.connect()
	}()
}

// newTCPServer will create a new TCP server using the address
func newTCPServer(addr string, requests chan<- *keystore.Request) *TCPServer {
	return &TCPServer{addr: addr, requests: requests, quit: make(chan bool)}
}

// connect will enable the event listener for incoming data
func (server *TCPServer) connect() {
	if server.connected {
		log.Println("The TCP server is already connected")
		return
	}

	// Make the connection
	var err error
	server.listener, err = net.Listen("tcp", server.addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("TCP server now connected to address: %s", server.addr)
	server.connected = true

	go func() {
		for {
			select {
			case <-server.quit:
				log.Println("TCP server connection is shutting down")

				// Then close the connection
				server.listener.Close()
				server.connected = false
				break
			default:
				// Wait for a connection on the listener
				conn, err := server.listener.Accept()
				if err != nil {
					log.Printf("An error occurred accepting a TCP connection: %s", err)
					continue
				}

				// Spin off a goroutine to handle this connection
				server.handleClient(conn)
			}
		}
	}()
}

// Disconnect will send the shutdown signal to this server connection
func (server *TCPServer) Disconnect() {

	// Spawn off the request to shutdown
	go func() {
		server.quit <- true
	}()
}

// handleClient will create a new client connector that will sit and listen
// for requests
func (server *TCPServer) handleClient(conn net.Conn) {

	// Create a new client connector
	newTCPClientHandler(conn, server.requests).start()
}

// CLIENT HANDLER

// newTCPClientHandler will wrap the client connection
// and listen for new requests
func newTCPClientHandler(conn net.Conn, requests chan<- *keystore.Request) *TCPClientHandler {
	client := &TCPClientHandler{requestChannel: requests, conn: conn, quit: make(chan bool)}
	client.encoder = gob.NewEncoder(conn)
	client.decoder = gob.NewDecoder(conn)
	return client
}

// Start will start the event listener for incoming data
func (tcp *TCPClientHandler) start() {
	clientaddr := tcp.conn.RemoteAddr().String()
	log.Printf("Received new client TCP connection: %s", clientaddr)

	// Listen for requests to send on the channel
	go func() {
		for {
			// Wait for the request from the client
			request := &keystore.Request{}
			err := tcp.decoder.Decode(request)
			if err != nil {
				// Either the data is incorrect or they have close the connection.
				// In both cases we shall also close the connection
				log.Printf("Client [%s] has closed the TCP connection", clientaddr)
				tcp.conn.Close()

				// Exit out of the routine
				return
			}
			log.Printf("Received TCP request from client: [%s]", clientaddr)

			// Wait for the response and send it back to the client
			go func() {

				// The channel can not be sent so will be created
				request.ResponseChannel = make(chan *keystore.Response)

				// Now send the request on the request channel
				tcp.requestChannel <- request
				response := <-request.ResponseChannel

				// We want to send the response back to the client
				log.Printf("Received response... Sending TCP response to client [%s]", clientaddr)
				err = tcp.encoder.Encode(response)
			}()
		}
	}()
}
