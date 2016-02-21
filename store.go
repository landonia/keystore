// Landon Wainwright.

// Package keystore provides an in memory key/value store service library
package keystore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

// Store creates a new in-memory store that can be read or written to disk
type Store struct {
	filePath string
	values   map[string]interface{}
}

// NewEmptyStore creates a new empty Store purely in memory and backed by no store
func NewEmptyStore() *Store {
	return NewStoreFromFile("")
}

// NewStoreFromFile creates a new empty Store that is backed by disk
func NewStoreFromFile(filePath string) *Store {
	return &Store{filePath: filePath, values: make(map[string]interface{})}
}

// UpdateFilePath will update the current file path to allow the data to be saved
// to that location
func (s *Store) UpdateFilePath(filePath string) {
	s.filePath = filePath
}

// ReadFromDisk loads the configuration from disk into this Store
// It returns an error if the data cannot be loaded from disk
func (s *Store) ReadFromDisk() (err error) {

	// If there is no file path then return
	if s.filePath != "" {
		log.Printf("Loading keystore from disk path %s", s.filePath)

		// Attempt to open the file
		f, err := os.Open(s.filePath)
		if err == nil {
			defer f.Close()
			var b bytes.Buffer
			_, err = b.ReadFrom(f)
			if err == nil {
				err = json.Unmarshal(b.Bytes(), &s.values)
			}
		}
	}
	return
}

// SaveToDisk will flush the current to disk if their is a valid filepath
func (s *Store) SaveToDisk() (err error) {

	// If there is no file path then return
	if s.filePath != "" {
		log.Printf("Saving keystore to disk path %s", s.filePath)

		// Marshall this to disk
		b, err := json.Marshal(s.values)
		if err == nil {

			// Get access to the file
			fo, err := os.Create(s.filePath)
			defer fo.Close()
			if err == nil {

				// Write the bytes to the buffer
				var buffer bytes.Buffer
				_, err = buffer.Write(b)
				if err == nil {

					// Write the buffer to disk
					_, err = buffer.WriteTo(fo)
				}
			}
		}
	}
	return
}

// generateError will return a new error containing the message
func generateError(message string) error {
	return errors.New(message)
}

// generateTypeError will return an error indicating that the value type for the key is incorrect
func generateTypeError(key string) error {
	return generateError(fmt.Sprintf("The value for key '%s' is not of the correct type", key))
}

// KeyExists returns true if the key exists in the store
func (s *Store) KeyExists(key string) (exists bool) {
	_, exists = s.values[key]
	return
}

// DeleteKey will delete the key from the store
func (s *Store) DeleteKey(key string) {
	delete(s.values, key)
}

// GetValue implements KeyValueStore interface
func (s *Store) GetValue(key string) (val interface{}, err error) {
	val, exists := s.values[key]

	// Check that the key exists and return an error if not
	if !exists {
		err = generateError(fmt.Sprintf("The key '%s' specified does not exist", key))
	}
	return
}

// GetBool implements KeyValueStore interface
func (s *Store) GetBool(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.(bool)
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// GetInt implements KeyValueStore interface
func (s *Store) GetInt(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.(int)
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// GetFloat implements KeyValueStore interface
func (s *Store) GetFloat(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.(float64)
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// GetString implements KeyValueStore interface
func (s *Store) GetString(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.(string)
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// GetArray implements KeyValueStore interface
func (s *Store) GetArray(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.([]interface{})
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// GetMap implements KeyValueStore interface
func (s *Store) GetMap(key string) (interface{}, error) {
	// Get the value and then attempt to type assert
	raw, err := s.GetValue(key)
	if err == nil {
		// Type assertion
		val, ok := raw.(map[string]interface{})
		if !ok {
			err = generateTypeError(key)
		}
		return val, err
	}
	return nil, err
}

// SetValue implements KeyValueStore interface
func (s *Store) SetValue(key string, value interface{}) error {
	s.values[key] = value
	return nil
}

// SetBool implements KeyValueStore interface
func (s *Store) SetBool(key string, value interface{}) error {
	val, ok := value.(bool)
	return s.setValueOrReturnError(key, val, ok)
}

// SetInt implements KeyValueStore interface
func (s *Store) SetInt(key string, value interface{}) error {
	val, ok := value.(int)
	return s.setValueOrReturnError(key, val, ok)
}

// SetFloat implements KeyValueStore interface
func (s *Store) SetFloat(key string, value interface{}) error {
	val, ok := value.(float64)
	return s.setValueOrReturnError(key, val, ok)
}

// SetString implements KeyValueStore interface
func (s *Store) SetString(key string, value interface{}) error {
	val, ok := value.(string)
	return s.setValueOrReturnError(key, val, ok)
}

// SetArray implements KeyValueStore interface
func (s *Store) SetArray(key string, value interface{}) error {
	val, ok := value.([]interface{})
	return s.setValueOrReturnError(key, val, ok)
}

// SetMap implements KeyValueStore interface
func (s *Store) SetMap(key string, value interface{}) error {
	val, ok := value.(map[string]interface{})
	return s.setValueOrReturnError(key, val, ok)
}

// setValueOrReturnError expects the value and whether the type assertion is ok.
// If the assertion is !ok an error is returned and the value is not set
func (s *Store) setValueOrReturnError(key string, val interface{}, ok bool) (err error) {
	if !ok {
		err = generateTypeError(key)
	} else {
		s.SetValue(key, val)
	}
	return
}
