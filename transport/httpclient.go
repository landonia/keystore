// Landon Wainwright.

package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/landonia/keystore"
)

// HTTPClient holds the HTTP client connection
type HTTPClient struct {
	*keystore.Sync           // Adopt the sync struct
	hostaddr       string    // the address to bind to
	quit           chan bool // The channel to wait on to finish the connection
	connected      bool      // Whether the server is currently connected
}

// NewHTTPClient will create a new HTTP connection using the host address
func NewHTTPClient(hostaddr string) *HTTPClient {
	return &HTTPClient{&keystore.Sync{RequestChannel: make(chan *keystore.Request)}, hostaddr, make(chan bool), false}
}

// Connect will start the event listener for incoming data
func (client *HTTPClient) Connect() {
	if client.connected {
		log.Printf("The HTTP client is already wrapped to host: %s", client.hostaddr)
		return
	}

	// Make the connection
	log.Printf("HTTP client wrapped to host: %s", client.hostaddr)
	client.connected = true

	// Listen for requests to send on the channel
	go func() {
		for {
			select {
			case request := <-client.RequestChannel:
				log.Println("Received a new client request")

				// Create the correct URL for the key
				var url string
				if client.hostaddr[0:7] != "http://" {
					url = fmt.Sprintf("http://%s/%s", client.hostaddr, request.Key)
				} else {
					url = fmt.Sprintf("%s/%s", client.hostaddr, request.Key)
				}

				// The HTTP request is based on the type of keystore operation
				var resp *http.Response
				var err error
				switch request.Op {
				case keystore.READ:

					// Make a GET request
					log.Printf("Making GET request: %s", url)
					resp, err = http.Get(url)
				case keystore.WRITE:

					// Encode the value to send in the body
					var b []byte
					if b, err = json.Marshal(request.Value.Val); err != nil {
						log.Printf("An error occurred marshalling GET request [%s] content: %s", url, err)
					} else {
						// Make a POST request
						log.Printf("Making POST request: %s", url)
						resp, err = http.Post(url, "application/json", bytes.NewBuffer(b))
					}
				case keystore.DELETE:
					// Make a DELETE request
					var req *http.Request
					if req, err = http.NewRequest("DELETE", url, nil); err != nil {
						log.Printf("An error occurred making DELETE request: %s", err)
					} else {
						// Make the request
						resp, err = http.DefaultClient.Do(req)
					}
				}

				// Check if there was an error
				if err != nil {
					log.Printf("An error occurred making the HTTP request [%s]: %s", url, err)
				} else {

					// Now just handle the response by gob'ling it up
					response := &keystore.Response{}
					if err = json.NewDecoder(io.LimitReader(resp.Body, MaxRequestLength)).Decode(response); err != nil {
						log.Printf("An error occurred decoding HTTP response: %s", err)
					}
					log.Println("Received response from HTTP request")

					// Send the response
					go func() { request.ResponseChannel <- response }()
				}
			case <-client.quit:
				log.Println("Client connection is shutting down")

				// Will shutdown the routine
				break
			}
		}
	}()
}

// Close will stop this client connection
func (client *HTTPClient) Close() {

	// Spawn off the request to shutdown
	go func() {
		client.quit <- true
	}()
}

// SendRequest will push the request onto the channel
func (client *HTTPClient) SendRequest(request *keystore.Request) {
	// Spawn off the request to the channel
	go func() {
		client.RequestChannel <- request
	}()
}
