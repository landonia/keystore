package main

import (
	"log"
	"time"

	"github.com/landonia/keystore/transport"
)

// Arb is an arbitrary data type that will be saved
type Arb struct {
	Myint    int
	Mystring string
	Myfloat  float32
}

// main will bootstrap the client libraries, make some requests and
// then check print out the result
func main() {
	// TCP

	go func() {
		tcpClient := transport.NewTCPClient("localhost:8081")
		tcpClient.Connect()

		// Make a write key request
		if err := tcpClient.SetString("tcpkey", "My tcpkey value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now read that key
		if val, err := tcpClient.GetString("tcpkey"); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println("The value: ", val)
		}

		// Arbitrary struct value
		arb := &Arb{200, "MyKeyValue", 12.5}
		if err := tcpClient.SetValue("tcparbkey", arb); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now attempt to read that value back into the arb
		arbVal := &Arb{}
		if err := tcpClient.GetValueType("tcparbkey", arbVal); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println(arbVal)
		}

		// Make another set value to delete the key after
		if err := tcpClient.SetString("tcpkey2", "My tcpkey2 value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Make a delete request
		tcpClient.DeleteKey("tcpkey2")
		tcpClient.Close()
	}()

	// UDP

	go func() {
		udpClient := transport.NewUDPClient("localhost:8082", "localhost:0")
		udpClient.Connect()

		// Make a write key request
		if err := udpClient.SetString("udpkey", "My udpkey value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now read that key
		if val, err := udpClient.GetString("udpkey"); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println("The value: ", val)
		}

		// Arbitrary struct value
		arb := &Arb{300, "MyKeyValue", 12.5}
		if err := udpClient.SetValue("udparbkey", arb); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now attempt to read that value back into the arb
		arbVal := &Arb{}
		if err := udpClient.GetValueType("udparbkey", arbVal); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println(arbVal)
		}

		// Make another set value to delete the key after
		if err := udpClient.SetString("udpkey2", "My udpkey2 value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Make a delete request
		udpClient.DeleteKey("udpkey2")
		udpClient.Close()
	}()

	// HTTP

	go func() {
		httpClient := transport.NewHTTPClient("localhost:8080")
		httpClient.Connect()

		// Make a write key request
		if err := httpClient.SetString("httpkey", "My httpkey value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now read that key
		if val, err := httpClient.GetString("httpkey"); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println("The value: ", val)
		}

		// Arbitrary struct value
		arb := &Arb{400, "MyKeyValue", 12.5}
		if err := httpClient.SetValue("httparbkey", arb); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Now attempt to read that value back into the arb
		arbVal := &Arb{}
		if err := httpClient.GetValueType("httparbkey", arbVal); err != nil {
			log.Println("An error occurred: ", err)
		} else {
			log.Println(arbVal)
		}

		// Make another set value to delete the key after
		if err := httpClient.SetString("httpkey2", "My httpkey2 value"); err != nil {
			log.Println("An error occurred: ", err)
		}

		// Make a delete request
		httpClient.DeleteKey("httpkey2")
		httpClient.Close()
	}()

	<-time.NewTimer(time.Second * 10).C
}
