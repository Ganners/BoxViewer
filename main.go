package main

import (
	"boxviewer/boxapi"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BoxStorage struct {
	Documents    map[string]*boxapi.DocumentObject
	Sessions     map[string]*boxapi.SessionObject
	fileLocation string
}

func (bs *BoxStorage) Save() error {
	// @TODO - Implement save functionality
	return nil
}

func (bs *BoxStorage) Load() error {
	// @TODO - Implement load functionality
	return nil
}

type BoxViewerServer struct {
	API     *boxapi.BoxApi
	Storage *BoxStorage
	Addr    string
	Port    string
}

// Starts an update loop, should:
// + Update any DocumentObject that is queued, such that it ends up ready
// + Update any expired SessionObjects
func (bvs *BoxViewerServer) updateLoop() {

	log.Println("Running update loop iteration\n")

	for key, document := range bvs.Storage.Documents {

		if document.Status != "done" {
			log.Println("Status is %s, attempting to update", document.Status)
			err, doc := bvs.API.UpdateDocument(document.Id)
			if err != nil {
				log.Println(err)
			}
			bvs.Storage.Documents[key] = doc
			bvs.Storage.Save()
		}
	}

	time.Sleep(10 * time.Second)
	go bvs.updateLoop()
}

// Processes a file, should do the following:
//
// + Check if there's a session
//   + If Yes: Show the viewer
//   + If No: Create a channel and run the process
//
// + Wait on a read channel for the document to be ready, once it is ready
//   then show the file (recursive?)
func (bvs *BoxViewerServer) processFile(filePath string) {

	session, sessionFound := bvs.Storage.Sessions[filePath]

	if sessionFound {
	}

	document, found := bvs.Storage.Documents[filePath]
	bvs.processFile(filePath)

	if !found {

		if err, document := bvs.API.MultipartUpload(filePath); err != nil {
			log.Println(err)
		} else {
			bvs.Storage.Documents[filePath] = document
			bvs.Storage.Save()
		}
	}
}

func (bvs *BoxViewerServer) viewHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	if filePaths, found := r.Form["url"]; found {

		filePath := filePaths[0]
		bvs.processFile(filePath)

		fmt.Fprintf(w, "<h1>Loading box view for: %s</h1>", filePath)
	}
}

func (bvs *BoxViewerServer) infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Information</h1>")
	fmt.Fprintf(w, "<h2>Document Objects</h2>")
	fmt.Fprintf(w, "<ul>")
	for _, document := range bvs.Storage.Documents {
		fmt.Fprintf(w, "<ul>")
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "Type", document.Type)
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "Id", document.Id)
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "Status", document.Status)
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "CreatedAt", document.CreatedAt)
		fmt.Fprintf(w, "</ul>")
	}
	fmt.Fprintf(w, "</ul>")

	fmt.Fprintf(w, "<h2>Session Objects</h2>")
	fmt.Fprintf(w, "<ul>")
	for _, session := range bvs.Storage.Sessions {
		fmt.Fprintf(w, "<ul>")
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "Type", session.Type)
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "Id", session.Id)
		fmt.Fprintf(w, "<li><strong>%s</strong>: %s", "ExpiresAt", session.ExpiresAt)
		fmt.Fprintf(w, "</ul>")
	}
	fmt.Fprintf(w, "</ul>")
}

func (bvs *BoxViewerServer) ListenAndServe() {
	http.HandleFunc("/info/", bvs.infoHandler)
	http.HandleFunc("/view/", bvs.viewHandler)
	http.ListenAndServe(bvs.Addr+":"+bvs.Port, nil)
}

func NewBoxViewerServer(addr string, port string, apiKey string, fileLocation string) *BoxViewerServer {

	storage := &BoxStorage{
		Documents: make(map[string]*boxapi.DocumentObject),
		Sessions:  make(map[string]*boxapi.SessionObject),
	}
	if err := storage.Load(); err != nil {
		log.Fatal("Error loading from storage")
	}

	api := boxapi.NewBoxApi(apiKey, fileLocation)

	bvs := &BoxViewerServer{
		API:     api,
		Storage: storage,
		Addr:    addr,
		Port:    port,
	}

	go bvs.updateLoop()

	return bvs
}

func main() {

	var (
		apiKey       = flag.String("key", "", "Your key for the Box Viewer API")
		fileLocation = flag.String("location", "", "Your destination to store downloaded files")
		serverAddr   = flag.String("addr", "0.0.0.0", "The address the server should run on")
		serverPort   = flag.String("port", "8080", "The port the server should run on")
	)

	flag.Parse()
	server := NewBoxViewerServer(*serverAddr, *serverPort, *apiKey, *fileLocation)
	server.ListenAndServe()
}
