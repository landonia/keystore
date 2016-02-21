// Landon Wainwright.

// Package keystore provides an in memory key/value store service library
package keystore

import (
	"fmt"
	"log"
)

// Service is the wrapper for the in-memory data store service
type Service struct {
	*Sync                // Adopt the sync struct
	store *Store         // The in-memory store
	quit  chan chan bool // Uses the channel as a signal to shutdown
}

// NewService will initialise a new keystore
// If filePath is not an empty string the contents  of the file will initialise
// the store and the in-memory values will be written to the store on shutdown
func NewService(filePath string) *Service {

	// Create a new instance of the key store
	return &Service{&Sync{make(chan *Request)}, NewStoreFromFile(filePath), make(chan chan bool)}
}

// Start will bootstrap the keystore service ready to receive requests
func (ks *Service) Start() {
	log.Println("Starting Keystore Service")
	ks.store.ReadFromDisk()

	// Spawn the store handler in a new go routine that will sit and wait for
	// operation requests. It is concurrently safe using channel blocking
	// for operations
	go func() {

		// Loop until it receives a message on the quite channel
		for {
			select {
			case request := <-ks.RequestChannel:
				response := &Response{}
				// A request has been made to perform an operation on the store
				switch request.Op {
				case READ:
					ks.readValue(request, response)
				case WRITE:
					ks.writeValue(request, response)
				case DELETE:
					ks.deleteKey(request, response)
				}

				// Send the response over the response channel
				go func() {
					request.ResponseChannel <- response
				}()
			case complete := <-ks.quit:
				// The signal to shutdown has been received

				// Write the values to disk
				if err := ks.store.SaveToDisk(); err != nil {
					log.Println(fmt.Errorf("Error saving values to disk: %s", err.Error()))
				}
				complete <- true
				return
			}
		}
	}()
}

// Stop will shutdown the keystore
func (ks *Service) Stop() chan bool {
	log.Println("Stopping Keystore Service")
	complete := make(chan bool)
	go func() {

		// Shutdown by ending the main store hander routine
		ks.quit <- complete
	}()
	return complete
}

// readValue will return the value from the store for the particular
// type and put that value into the Response
func (ks *Service) readValue(request *Request, response *Response) {

	// Create a new value holder for this response
	response.Value = &ValueHolder{Type: request.Value.Type}
	var err error

	// The get method is called for the specific type requested which
	// performs type assertion here to ensure that the correct type
	// is received.
	switch request.Value.Type {
	case BOOL:
		response.Value.Val, err = ks.store.GetBool(request.Key)
	case INT:
		response.Value.Val, err = ks.store.GetInt(request.Key)
	case FLOAT:
		response.Value.Val, err = ks.store.GetFloat(request.Key)
	case STRING:
		response.Value.Val, err = ks.store.GetString(request.Key)
	case ARRAY:
		response.Value.Val, err = ks.store.GetArray(request.Key)
	case MAP:
		response.Value.Val, err = ks.store.GetMap(request.Key)
	default:
		response.Value.Val, err = ks.store.GetValue(request.Key)
	}

	// if no error occurred during this operation then the request was a success
	if err == nil {
		response.Success = true
	} else {
		response.Success = false
		response.Error = err.Error()
	}
}

// writeValue will write a value to the store
func (ks *Service) writeValue(request *Request, response *Response) {
	var err error

	// The value passed to all methods is of type interface{} but they each
	// check that the value is of the correct type. This is done here as it will
	// ensure that the specific value type is correct without this having to be
	// done on each end of the request and response. The receiver will only have to
	// assert the type when they are receiving NONE particular type
	switch request.Value.Type {
	case BOOL:
		err = ks.store.SetBool(request.Key, request.Value.Val)
	case INT:
		err = ks.store.SetInt(request.Key, request.Value.Val)
	case FLOAT:
		err = ks.store.SetFloat(request.Key, request.Value.Val)
	case STRING:
		err = ks.store.SetString(request.Key, request.Value.Val)
	case ARRAY:
		err = ks.store.SetArray(request.Key, request.Value.Val)
	case MAP:
		err = ks.store.SetMap(request.Key, request.Value.Val)
	default:
		ks.store.SetValue(request.Key, request.Value.Val)
	}

	// if no error occurred during this operation then the request was a success
	if err == nil {
		response.Success = true
	} else {
		response.Success = false
		response.Error = err.Error()
	}
}

// deleteKey will delete the key and value from the store
func (ks *Service) deleteKey(request *Request, response *Response) {

	// Then delete the key if it is present
	ks.store.DeleteKey(request.Key)
	response.Success = true
}
