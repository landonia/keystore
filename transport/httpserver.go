// Landon Wainwright.

package transport

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/landonia/keystore"
)

// MaxRequestLength specifies the amount of bytes accepted on a request
// Do not let someone hang the service by sending continous stream of data
// on the request. 10Kb will be big enough for this example.
const MaxRequestLength int64 = 1024

// StartHTTPServer will start a new HTTP server allowing requests
// to be made to the key store service over a REST interface
func StartHTTPServer(addr string, requestChannel chan<- *keystore.Request) {
	http.Handle("/", generateHandler(requestChannel, operationHandler))

	// Start the server
	go func() {
		log.Printf("Starting HTTP server using address: %s", addr)

		// Process will stop when the program has finished or an error occurs
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}

// generateHandler will generate a handler that closes over the state required
func generateHandler(requestChannel chan<- *keystore.Request, handler func(http.ResponseWriter, *http.Request, chan<- *keystore.Request)) http.Handler {

	// Handler that passes the requests along with the request and response
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler(w, r, requestChannel) })
}

// operationHandler will return the correct status messages for any requests made
// using incorrect paths
func operationHandler(w http.ResponseWriter, r *http.Request, requestChannel chan<- *keystore.Request) {
	var request *keystore.Request
	// Extract the key name from the URL
	if key := r.URL.Path[1:]; key == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if r.Method == "POST" {
		// Read the body into a string for json decoding
		var content interface{}
		err := json.NewDecoder(io.LimitReader(r.Body, MaxRequestLength)).Decode(&content)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// The key and value are now retrieved
		request = keystore.NewWriteRequest(key, keystore.NONE, content)
	} else if r.Method == "GET" {

		// A read request
		request = keystore.NewReadRequest(key, keystore.NONE)
	} else if r.Method == "DELETE" {
		// A delete request
		request = keystore.NewDeleteRequest(key)
	}

	// Send the request to the keystore
	requestChannel <- request

	// Get the response channel from the request
	select {
	case response := <-request.ResponseChannel:
		// Marshall the response
		content, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the correct content type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		// Write the content back
		w.Write(content)
	case <-time.NewTimer(time.Minute * 1).C:

		// Wait for a minute and then send a timeout
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
