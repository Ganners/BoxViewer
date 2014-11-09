package main

import (
	"flag"
	"log"
	"net/http"
)

// Main function is responsible for parsing our flags and starting
// the server, as well as initialising the workers
func main() {

	// Parse all of our flags for our application
	var (
		serverAddr   = flag.String("addr", "0.0.0.0", "The address the server should run on")
		serverPort   = flag.String("port", "8080", "The port the server should run on")
		apiKey       = flag.String("key", "", "Your key for the Box Viewer API")
		fileLocation = flag.String("location", "", "Your destination to store downloaded files")
	)
	flag.Parse()

	// Produce a worker pool as our jobs take a long time, creates
	// a channel to store jobs in progress and jobs complete
	bw := NewBoxWorker(*apiKey, *fileLocation)
	jobs, results := bw.WorkerPool(5)

	// Create a RequestMux which returns a requests channel, this
	// directs all of our stuff to avoid duplication
	requests := RequestMux(jobs, results)

	// Produce and start our server, passing in the requests channel
	server := &Server{Requests: requests}

	// Serve and do magic
	log.Fatal(http.ListenAndServe(
		*serverAddr+":"+*serverPort, server))
}
