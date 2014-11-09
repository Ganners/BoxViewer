package main

import (
	"log"
	"net/http"
)

// A job consists of a filename (input) and a result (output)
type Job struct {
	Filename string
	Result   string
}

// Request has a Job (Filename and Result string) and a channel
// of strings for the result
type Request struct {
	Job        *Job
	ResultChan chan string
}

// Our RequestMux bundles together our jobs into a channel of requests
// (which is returned). Similar job filenames get queued together
func RequestMux(jobs chan *Job, results chan *Job) chan *Request {
	requests := make(chan *Request)

	go func() {
		queues := make(map[string][]*Request)

		for {
			select {
			case request := <-requests:
				job := request.Job

				// Append to the queue filename slice
				queues[job.Filename] = append(queues[job.Filename], request)

				// If the length of the queue is one then pass to the jobs channel
				if len(queues[job.Filename]) == 1 {
					go func() {
						jobs <- job
					}()
				}

			case job := <-results:
				// If we receive a result then pass to the results channel on
				// the job (and remove it from the queue)
				for _, request := range queues[job.Filename] {
					request.ResultChan <- job.Result
				}

				delete(queues, job.Filename)
			}
		}
	}()

	return requests
}

type Server struct {
	Requests chan *Request
}

// Our HTTP server, will listen for a HTTP request, creating a new
// request object and pass that into our Mux. The Mux then bundles
// it together and sends it to our WorkerPool -> Worker
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Listen for a filename
	filename := req.URL.Query().Get("filename")

	// Produce our request with a result channel
	request := &Request{
		Job:        &Job{Filename: filename},
		ResultChan: make(chan string),
	}

	// Pass all requests to our RequestMux
	s.Requests <- request
	path := <-request.ResultChan // Block
	log.Println("Redirecting to", path)
	http.Redirect(w, req, path, 302)
}
