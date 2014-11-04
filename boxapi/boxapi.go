package boxapi

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type BoxApi struct {
	ApiKey       string
	MultipartUrl string
	DocumentUrl  string
	SessionUrl   string
}

func (box *BoxApi) generateUniqueFilename(filePath string) string {
	hash := md5.Sum([]byte(filePath))
	return hex.EncodeToString(hash[:])
}

func (box *BoxApi) downloadFile(filePath string) (error, *os.File) {

	// Check if the file that has been passed exists and is reachable
	resp, err := http.Get(filePath)
	if err != nil {
		return err, nil
	}

	fmt.Printf("%v", resp)

	return nil, nil
}

func (box *BoxApi) MultipartUpload(filePath string) (err error, docObj *DocumentObject) {

	fmt.Printf("Hi all")

	// Get file
	err, file := box.downloadFile(filePath)
	if err != nil {
		return err, nil
	}

	if err != nil {
		return err, nil
	}

	return nil, nil

	// Perform the upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "document")
	if err != nil {
		return err, nil
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err, nil
	}

	err = writer.Close()
	if err != nil {
		return err, nil
	}

	mprequest, err := http.NewRequest("POST", box.MultipartUrl, body)
	if err != nil {
		return err, nil
	}

	client := &http.Client{}
	mpresponse, err := client.Do(mprequest)
	if err != nil {
		return err, nil
	}
	defer mpresponse.Body.Close()

	mybody, err := ioutil.ReadAll(mpresponse.Body)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(mybody, &docObj)
	if err != nil {
		return err, nil
	}

	return nil, docObj
}

func (box *BoxApi) GetDocument(documentId string) *DocumentObject {
	return &DocumentObject{}
}

func (box *BoxApi) GetSession(documentId string) *SessionObject {
	return &SessionObject{}
}

func NewBoxApi(key string) *BoxApi {
	return &BoxApi{
		ApiKey:       key,
		MultipartUrl: "https://upload.view-api.box.com",
		DocumentUrl:  "https://view-api.box.com/1/documents",
		SessionUrl:   "https://view-api.box.com/1/sessions",
	}
}
