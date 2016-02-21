// Landon Wainwright.

package transport

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/landonia/keystore"
)

// UDPClient holds the UDP client connection
type UDPClient struct {
	*keystore.Sync              // Adopt the sync struct
	hostaddr       string       // The host address to connect to
	localaddr      string       // The local address to bind to
	conn           *net.UDPConn // The udp connection
	quit           chan bool    // The channel to wait on to finish the connection
	connected      bool         // Whether the server is currently connected
}

// NewUDPClient will create a new UDP connection using the host address
func NewUDPClient(hostaddr, localaddr string) *UDPClient {
	return &UDPClient{&keystore.Sync{RequestChannel: make(chan *keystore.Request)}, hostaddr, localaddr, nil, make(chan bool), false}
}

// Connect will start the event listener for incoming data
func (udp *UDPClient) Connect() {
	if udp.connected {
		log.Println("The UDP client is already connected")
		return
	}

	// Get the local address to receive the response
	localaddr, _ := net.ResolveUDPAddr("udp", udp.localaddr)

	// Get the server address to send the data
	serveraddr, _ := net.ResolveUDPAddr("udp", udp.hostaddr)

	// Make the connection
	var err error
	udp.conn, err = net.DialUDP("udp", localaddr, serveraddr)
	if err != nil {
		log.Fatal("An error occurred whilst making the UDP connection: ", err)
		return
	}
	log.Printf("UDP client now connected to remote address: %s: local address: %s", udp.hostaddr, udp.localaddr)
	udp.connected = true

	// Listen for requests to send on the channel
	go func() {
		for {
			select {
			case request := <-udp.RequestChannel:
				log.Println("Received a new client request")

				// Use gob to get a stream of bytes to write to a packet
				var buff bytes.Buffer
				if err := gob.NewEncoder(&buff).Encode(request); err != nil {
					log.Printf("Error writing request to buffer: %s", err)
				}

				// Write the request bytes to the client
				_, err := udp.conn.Write(buff.Bytes())
				if err != nil {
					log.Printf("Error writing request to client: %s", err)
				}
				log.Println("Waiting for client response")

				// Wait for the response
				buf := make([]byte, 1024)
				n, err := udp.conn.Read(buf)
				if err != nil {
					log.Println("Error whilst reading UDP response packet: ", err)
				}
				response := &keystore.Response{}
				if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(response); err != nil {
					log.Printf("Error whilst reading UDP packet: %s", err)
					if response.Error == "" {
						response.Error = err.Error()
					}
				}
				log.Println("Received response from server")

				// Send the response
				go func() { request.ResponseChannel <- response }()
			case <-udp.quit:
				log.Println("Client connection is shutting down")

				// Close the connection
				udp.conn.Close()
				break
			}
		}
	}()
}

// Close will stop this client connection
func (udp *UDPClient) Close() {

	// Spawn off the request to shutdown
	go func() {
		udp.quit <- true
	}()
}

// SendRequest will push the request onto the channel
func (udp *UDPClient) SendRequest(request *keystore.Request) {
	// Spawn off the request to the channel
	go func() {
		udp.RequestChannel <- request
	}()
}
