package main

import (
	"boxviewer/boxapi"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type BoxStorage struct {
	Documents    map[string]*boxapi.DocumentObject
	Sessions     map[string]*boxapi.SessionObject
	fileLocation string
}

func (bs *BoxStorage) save() error {
	// @TODO - Implement save functionality
	return nil
}

func (bs *BoxStorage) load() error {
	// @TODO - Implement load functionality
	return nil
}

type BoxViewerServer struct {
	API     *boxapi.BoxApi
	Storage *BoxStorage
	Addr    string
	Port    string
}

func (bvs *BoxViewerServer) viewHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if filePaths, found := r.Form["url"]; found {

		filePath := filePaths[0]
		value, found := bvs.Storage.Documents[filePath]

		if !found {
			if err, document := bvs.API.MultipartUpload(filePath); err != nil {
				log.Println(err)
			} else {
				bvs.Storage.Documents[filePath] = document
				value = document
			}
		}

		fmt.Fprintf(w, "<h1>Loading box view for: %s</h1>", filePath)
		fmt.Fprintf(w, "<h1>Value is: %v</h1>", value)
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
	if err := storage.load(); err != nil {
		log.Fatal("Error loading from storage")
	}

	api := boxapi.NewBoxApi(apiKey, fileLocation)

	return &BoxViewerServer{
		API:     api,
		Storage: storage,
		Addr:    addr,
		Port:    port,
	}
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

	// api := boxapi.NewBoxApi(*apiKey, *fileLocation)
	// err, docObj := api.MultipartUpload(
	// 	"http://www.mmta.co.uk/uploads/2014/09/26/102751_crucible_sept_14_final_v2.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%v %v", err, docObj)
}
