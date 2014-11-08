package main

import (
	"boxviewer/boxapi"
	"flag"
	"fmt"
	"log"
)

var (
	apiKey       = flag.String("key", "", "Your key for the Box Viewer API")
	fileLocation = flag.String("location", "", "Your destination to store downloaded files")
)

func main() {
	flag.Parse()
	api := boxapi.NewBoxApi(*apiKey, *fileLocation)

	err, docObj := api.MultipartUpload(
		"http://www.mmta.co.uk/uploads/2014/09/26/102751_crucible_sept_14_final_v2.pdf")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v %v", err, docObj)
}
