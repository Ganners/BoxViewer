package main

import "boxviewer/boxapi"

func main() {
	api := boxapi.NewBoxApi("12345")
	api.MultipartUpload(
		"http://www.mmta.co.uk/uploads/2014/09/26/102751_crucible_sept_14_final_v2.pdf")
}
