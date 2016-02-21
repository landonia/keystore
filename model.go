// Landon Wainwright.

// Package keystore provides an in memory key/value store service library
package keystore

// Op is the operation type for the request to the data store
type Op uint

// Flags for the request operation type
const (
	READ   Op = 1 << iota // A request to read the value
	WRITE  Op = 1 << iota // A request to write a value
	DELETE Op = 1 << iota // A request to delete the key and value
)

// Type allows the requester to specify the type of data it is expecting
// This allows the caller to either handle the value or the error directly without
// having to check the type. NONE can be specified meaning any type will be accepted
type Type uint

// Flags the data types
const (
	BOOL   Type = 1 << iota // Expecting a boolean
	INT    Type = 1 << iota // Expecting an int
	FLOAT  Type = 1 << iota // Expecting a float
	STRING Type = 1 << iota // Expecting a string
	ARRAY  Type = 1 << iota // Expecting an array
	MAP    Type = 1 << iota // Expecting a map
	NONE   Type = 1 << iota // Expecting any type
)

// Request are the requests that will be sent over the channel
// for operations on the store, such as a read or write
type Request struct {
	Op              Op             // The operation required
	Key             string         // The key (used for all requests)
	Value           *ValueHolder   // The request value (used for write requests only)
	ResponseChannel chan *Response // The return channel
}

// Response will be returned for each Request containing the result of the operation.
// If the data type had been specified the value will be contained within the
// specific value type objects otherwise it will be placed in the arbitrary
// value field. On a write operation the value returned will be the existing value.
type Response struct {
	Success bool         // True if the operation was a success (Error may still be present for write and delete operations)
	Error   string       // Will contain any errors
	Value   *ValueHolder // The response values
}

// ValueHolder wraps the value, but if no error the value will be of the type expected
type ValueHolder struct {
	Type Type        // The data type expected (if read Op)
	Val  interface{} // The arbitrary value if NONE DataType specified
}

// KeyValueStore Defines the interface for a simple key value store type.
// It can be used by both the server and the clients to provide a uniform
// API that only differs by the transport method
type KeyValueStore interface {
	// GetValueType This is used to attempt to parse the value into the interface value provided
	// See https://golang.org/pkg/encoding/json/#Unmarshal for more information
	GetValueType(key string, v interface{}) error

	// GetValue will return the raw value that has been stored
	GetValue(key string) (interface{}, error)

	// GetBool returns a bool for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetBool(key string) (interface{}, error)

	// GetInt returns an int for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetInt(key string) (interface{}, error)

	// GetFloat returns an float for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetFloat(key string) (interface{}, error)

	// GetString returns a string for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetString(key string) (interface{}, error)

	// GetArray returns an array type for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetArray(key string) (interface{}, error)

	// GetMap returns a map type for the key specified or if the key does not exist
	// or the value is not of the correct type an error is returned
	GetMap(key string) (interface{}, error)

	// SetValue will store an arbitrary type value for the key specified
	SetValue(key string, value interface{}) error

	// SetBool will store an boolean type value for the key specified
	SetBool(key string, value interface{}) error

	// SetInt will store an int type value for the key specified
	SetInt(key string, value interface{}) error

	// SetFloat will store a float type value for the key specified
	SetFloat(key string, value interface{}) error

	// SetString will store a string type value for the key specified
	SetString(key string, value interface{}) error

	// SetArray will store an array type value for the key specified
	SetArray(key string, value interface{}) error

	// SetMap will store a map type value for the key specified
	SetMap(key string, value interface{}) error

	// DeleteKey will delete the key and value from the store
	DeleteKey(key string)
}

// NewReadRequest will generate a new Request for reading a key
func NewReadRequest(key string, dType Type) *Request {
	return &Request{Op: READ, Key: key, Value: &ValueHolder{Type: dType}, ResponseChannel: make(chan *Response)}
}

// NewWriteRequest will generate a new Request for reading a key
func NewWriteRequest(key string, dType Type, value interface{}) *Request {
	return &Request{Op: WRITE, Key: key, Value: &ValueHolder{Type: dType, Val: value}, ResponseChannel: make(chan *Response)}
}

// NewDeleteRequest will generate a new Request for reading a key
func NewDeleteRequest(key string) *Request {
	return &Request{Op: DELETE, Key: key, Value: &ValueHolder{Type: NONE}, ResponseChannel: make(chan *Response)}
}
