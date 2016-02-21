// Landon Wainwright.

package transport

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/landonia/keystore"
)

// TCPClient holds the TCP client connection
type TCPClient struct {
	*keystore.Sync              // Adopt the sync struct
	hostaddr       string       // the address to bind to
	conn           net.Conn     // The tcp connection
	quit           chan bool    // The channel to wait on to finish the connection
	encoder        *gob.Encoder // The encoder for this connection
	decoder        *gob.Decoder // The decoder for this connection
	connected      bool         // Whether the server is currently connected
}

// NewTCPClient will create a new TCP connection using the host address
func NewTCPClient(hostaddr string) *TCPClient {
	return &TCPClient{&keystore.Sync{RequestChannel: make(chan *keystore.Request)}, hostaddr, nil, make(chan bool), nil, nil, false}
}

// Connect will start the event listener for incoming data
func (client *TCPClient) Connect() {
	if client.connected {
		log.Println("The TCP client is already connected")
		return
	}

	// Make the connection
	var err error
	client.conn, err = net.Dial("tcp", client.hostaddr)
	if err != nil {
		log.Fatal("An error occurred whilst making the connection: ", err)
		return
	}
	log.Printf("TCP client now connected to address: %s", client.hostaddr)
	client.encoder = gob.NewEncoder(client.conn)
	client.decoder = gob.NewDecoder(client.conn)
	client.connected = true

	// Listen for requests to send on the channel
	go func() {
		for {
			select {
			case request := <-client.RequestChannel:
				log.Println("Received a new client request")

				// Use the encoder to send the request directly
				// NOTE - It would make sense in reality to use protocol buffers here to allow
				// other systems to easily encode/decode the required payload
				if err := client.encoder.Encode(request); err != nil {
					log.Printf("An error occurred encoding TCP request: %s", err)
				}
				log.Println("Waiting for client response")

				// Wait for the response
				response := &keystore.Response{}
				if err := client.decoder.Decode(response); err != nil {
					log.Printf("An error occurred decoding TCP request: %s", err)
					if response.Error == "" {
						response.Error = err.Error()
					}
				}
				log.Println("Received response from server")

				// Send the response
				go func() { request.ResponseChannel <- response }()
			case <-client.quit:
				log.Println("Client connection is shutting down")

				// Close the connection
				client.conn.Close()
				break
			}
		}
	}()
}

// Close will stop this client connection
func (client *TCPClient) Close() {

	// Spawn off the request to shutdown
	go func() {
		client.quit <- true
	}()
}

// SendRequest will push the request onto the channel
func (client *TCPClient) SendRequest(request *keystore.Request) {
	// Spawn off the request to the channel
	go func() {
		client.RequestChannel <- request
	}()
}
