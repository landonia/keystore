# keystore

An in-memory key/value store written using the Go Language.

## Overview

The keystore can be loaded from disk and will be flushed to disk on shutdown.
You have the option of including the keystore as a library in an
existing project or it can be run independently as a service by executing
cmd/keystore. You can choose from the transport layers provided including
TCP/UDP/HTTP. The clients are provided to access the store.

## Maturity

This is the first stab. I need to add some better fine grained error handling.
I will add a request type that will allow you to flush the store to disk.

## Installation

For the library simply run `go get github.com/landonia/keystore`

For the service simply run `go get github.com/landonia/keystore/cmd/keystore`

## Execute Service

Using the provided command to run the keystore as an independent service.

  `keystore -dataPath ~/keystore/backup -httpAddr 8080 -tcpAddr 8081 -udpAddr 8082`

## Use as Library

	package main
	
	import (
  		"github.com/landonia/keystore"
  		"github.com/landonia/keystore/transport"
  	)
  	
  	func main() {
  	
  		// Create a key store in disk
  		keystore := keystore.NewService("path/to/file/location")

  		// Bind the network protocols that are required
  		transport.StartHTTPServer(":8080", keystore.RequestChannel)
  		transport.StartTCPServer(":8081", keystore.RequestChannel)
  		transport.StartUDPServer(":8082", keystore.RequestChannel)

  		// Start
  		keystore.Start()

  		// .... wait until application ends and then shutdown and wait
  		<-keystore.Stop()
  	}

## Client Library Example

You will find a simple example of using the client transport libraries in cmd/clientexample/main.go

## Future

This can be used now but I want to add more request types and much better fine grained
error handling (as it will assume any client error is a disconnect type right now).
I would also want to add protocol buffer data types to make it easier to create clients
in other languages.

## About

keystore was written by [Landon Wainwright](http://www.landotube.com) | [GitHub](https://github.com/landonia).

Follow me on [Twitter @landoman](http://www.twitter.com/landoman)!
