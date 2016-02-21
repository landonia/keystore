// Landon Wainwright.

// Package keystore provides an in memory key/value store service library
package keystore

import (
	"encoding/json"
	"errors"
)

// Sync implements the KeyValueStore and provides a synchronous blocking API
// that can be used with the keystore, server and client transports
type Sync struct {

	// All requests to the store are sent over this channel
	RequestChannel chan *Request
}

// waitForReadValue will block until the value has arrived
func waitForReadValue(requestChannel chan *Request, request *Request) (interface{}, error) {
	requestChannel <- request
	response := <-request.ResponseChannel
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Value.Val, nil
}

// GetValueType implements KeyValueStore
func (s *Sync) GetValueType(key string, v interface{}) error {

	// Using this impl the set value converts the item to a json string
	raw, err := s.GetString(key)
	if err != nil {
		return err
	}

	// Attempt to parse into the expected interface type
	return json.Unmarshal([]byte(raw.(string)), v)
}

// GetValue implements KeyValueStore
func (s *Sync) GetValue(key string) (interface{}, error) {

	// Using this impl the set value converts the item to a json string
	return s.GetString(key)
}

// GetBool implements KeyValueStore
func (s *Sync) GetBool(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, BOOL))
}

// GetInt implements KeyValueStore
func (s *Sync) GetInt(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, INT))
}

// GetFloat implements KeyValueStore
func (s *Sync) GetFloat(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, FLOAT))
}

// GetString implements KeyValueStore
func (s *Sync) GetString(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, STRING))
}

// GetArray implements KeyValueStore
func (s *Sync) GetArray(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, ARRAY))
}

// GetMap implements KeyValueStore
func (s *Sync) GetMap(key string) (interface{}, error) {
	return waitForReadValue(s.RequestChannel, NewReadRequest(key, MAP))
}

// waitForWriteValue will block until the response has arrived
func waitForWriteValue(requestChannel chan *Request, request *Request) error {
	requestChannel <- request
	response := <-request.ResponseChannel
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

// SetValue implements KeyValueStore
func (s *Sync) SetValue(key string, value interface{}) error {

	// So as to not have to register structs with gob and to handle arbitrary data
	// we will marhsall this to json (again, as mentioned in the servers/clients
	// I would usually use some intermediate protocol buffers)
	jsonb, _ := json.Marshal(value)
	return s.SetString(key, string(jsonb))
}

// SetBool implements KeyValueStore
func (s *Sync) SetBool(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, BOOL, value))
}

// SetInt implements KeyValueStore
func (s *Sync) SetInt(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, INT, value))
}

// SetFloat implements KeyValueStore
func (s *Sync) SetFloat(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, FLOAT, value))
}

// SetString implements KeyValueStore
func (s *Sync) SetString(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, STRING, value))
}

// SetArray implements KeyValueStore
func (s *Sync) SetArray(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, ARRAY, value))
}

// SetMap implements KeyValueStore
func (s *Sync) SetMap(key string, value interface{}) error {
	return waitForWriteValue(s.RequestChannel, NewWriteRequest(key, MAP, value))
}

// DeleteKey implements KeyValueStore
func (s *Sync) DeleteKey(key string) {
	waitForWriteValue(s.RequestChannel, NewDeleteRequest(key))
}
