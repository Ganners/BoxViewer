package main

import (
	"boxviewer/boxapi"
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"time"
)

const BoxStorageFileName = "box_storage_cache.bin"

// Our storage object, allows us to store documents and sessions in
// memory for quick retrieval. Should also allow caching to disk
type BoxStorage struct {
	Documents    map[string]*boxapi.DocumentObject
	Sessions     map[string]*boxapi.SessionObject
	fileLocation string
}

// Saves our BoxStorage to disk, will store in the fileLocation that
// the other files are also going to save into
func (bs *BoxStorage) Save() error {

	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(bs)

	if err != nil {
		return err
	}

	log.Println("Saving file to", bs.fileLocation+"/"+BoxStorageFileName)

	fh, err := os.OpenFile(bs.fileLocation+"/"+BoxStorageFileName,
		os.O_CREATE|os.O_WRONLY, 0777)
	defer fh.Close()

	if err != nil {
		return err
	}

	_, err = fh.Write(b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Reads our BoxStorage from disk (should it exist)
func (bs *BoxStorage) Load() error {

	log.Println("Loading file from", bs.fileLocation+"/"+BoxStorageFileName)

	fh, err := os.Open(bs.fileLocation + "/" + BoxStorageFileName)
	defer fh.Close()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(fh)
	err = dec.Decode(&bs)
	if err != nil {
		return err
	}

	return nil
}

// BoxWorker has a constructor so dependencies on the storage and
// on the API can be easily managed
type BoxWorker struct {
	API     *boxapi.BoxApi
	Storage *BoxStorage
}

// Our worker - does the tough job of fetching the URL of the box
// viewer that we'll redirect to
func (bw *BoxWorker) getBoxViewerURL(fileName string) string {

	// Check if there is a session that exists already that has not
	// expired. If yes create response
	if session, found := bw.Storage.Sessions[fileName]; found {
		if !session.IsExpired() {
			log.Println("Session has been found and is not expired for ", fileName)
			return bw.API.GetViewerURL(session.Id)
		}
	}

	// If there is no session or it has expired then check if there
	// is a document that exists.
	if document, found := bw.Storage.Documents[fileName]; found {

		// Check if the document is ready, if it isn't then we want to keep
		// trying it until it is complete
		for {
			if document.Status == "done" {
				break
			}

			// Time interval between attempts
			time.Sleep(1 * time.Second)
			_, document = bw.API.GetDocument(document.Id)
		}

		// If there is a document that exists, ask for a new session and
		// create the URL
		log.Println("Document found, generating new session for ", fileName)
		err, session := bw.API.GetSession(document.Id)

		if err != nil {
			log.Fatal(err)
		}

		bw.Storage.Sessions[fileName] = session
		bw.Storage.Save()
		return bw.API.GetViewerURL(session.Id)
	}

	// If there is not a document that exists, it is new and so we must
	// upload it
	log.Println("No document found, creating one for ", fileName)
	err, document := bw.API.MultipartUpload(fileName)
	if err != nil {
		log.Fatal(err)
	}
	bw.Storage.Documents[fileName] = document
	bw.Storage.Save()

	// May as well call ourself so we aren't rewriting lines of code
	return bw.getBoxViewerURL(fileName)
}

// The Worker is run as a goroutine with the jobs and results. It
// listens for jobs, when it receives one it does some hard work and
// then passes the result through to the results channel
func (bw *BoxWorker) Worker(jobs chan *Job, results chan *Job) {

	for job := range jobs {
		job.Result = bw.getBoxViewerURL(job.Filename)
		results <- job
	}
}

// Constructs a new BoxWorker with working API and storage
func NewBoxWorker(apiKey string, fileLocation string) *BoxWorker {

	bs := &BoxStorage{
		Documents:    make(map[string]*boxapi.DocumentObject),
		Sessions:     make(map[string]*boxapi.SessionObject),
		fileLocation: fileLocation,
	}

	// Try and load the cache from disk
	bs.Load()

	bw := &BoxWorker{
		API:     boxapi.NewBoxApi(apiKey, fileLocation),
		Storage: bs,
	}
	return bw
}

// Takes an integer for the number of workers in the pool. Creates
// a channel for jobs and results which are returned). For the number
// you put in we generate that many Worker goroutines
func (bw *BoxWorker) WorkerPool(n int) (chan *Job, chan *Job) {
	jobs := make(chan *Job)
	results := make(chan *Job)

	for i := 0; i < n; i++ {
		go bw.Worker(jobs, results)
	}

	return jobs, results
}
