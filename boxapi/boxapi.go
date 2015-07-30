package boxapi

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"errors"
)

type BoxApi struct {
	ApiKey       string
	FileLocation string
	MultipartUrl string
	DocumentUrl  string
	SessionUrl   string
}

type LocalFile struct {
	FileName string
	File     *os.File
}

func (box *BoxApi) generateUniqueFilename(filePath string) string {

	// Grab the extension to add on to the end of the md5
	ext := filepath.Ext(filePath)

	// Md5 and append the extension on
	hash := md5.Sum([]byte(filePath))
	return hex.EncodeToString(hash[:]) + ext
}

func (box *BoxApi) downloadFile(filePath string) (*LocalFile, error) {

	// Generate the filename from the filepath
	localFilename := box.FileLocation + "/" + box.generateUniqueFilename(filePath)

	if _, err := os.Stat(localFilename); os.IsNotExist(err) {
		// Check if the file that has been passed exists and is reachable
		resp, err := http.Get(filePath)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Write the file to disk with a unique filename
		err = ioutil.WriteFile(
			localFilename,
			contents,
			0777)

		if err != nil {
			return nil, err
		}
	}

	// Grab an os.File of the local file
	localFile, err := os.Open(localFilename)
	if err != nil {
		return nil, err
	}

	return &LocalFile{localFilename, localFile}, err
}

func (box *BoxApi) MultipartUpload(filePath string) (err error, docObj *DocumentObject) {

	// Get file (os.File)
	file, err := box.downloadFile(filePath)
	if err != nil {
		return err, nil
	}

	// Create a new form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.FileName))
	if err != nil {
		return err, nil
	}

	// Copy file contents into our form file for upload
	_, err = io.Copy(part, file.File)
	if err != nil {
		return err, nil
	}
	err = writer.Close()

	// Create our new request
	mprequest, err := http.NewRequest(
		"POST", box.MultipartUrl, body)
	mprequest.Header.Set("Authorization", "Token "+box.ApiKey)
	mprequest.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())

	if err != nil {
		return err, nil
	}

	// Perform the request, grab response and close
	client := &http.Client{}
	mpresponse, err := client.Do(mprequest)
	if err != nil {
		return err, nil
	}
	defer mpresponse.Body.Close()

	// Read the response into a byte slice
	mybody, err := ioutil.ReadAll(mpresponse.Body)
	if err != nil {
		return err, nil
	}

	// Unmarshal it into our document object
	err = json.Unmarshal(mybody, &docObj)
	if err != nil {
		return err, nil
	}

	if docObj.Type == "error" {
		return errors.New("Document could not be created"), nil
	}

	return nil, docObj
}

func (box *BoxApi) GetDocument(documentId string) (err error, docObj *DocumentObject) {

	// Create our new request
	mprequest, err := http.NewRequest(
		"GET", box.DocumentUrl+"/"+documentId, nil)
	mprequest.Header.Set("Authorization", "Token "+box.ApiKey)

	if err != nil {
		return err, nil
	}

	// Perform the request, grab response and close
	client := &http.Client{}
	mpresponse, err := client.Do(mprequest)
	if err != nil {
		return err, nil
	}
	defer mpresponse.Body.Close()

	// Read the response into a byte slice
	mybody, err := ioutil.ReadAll(mpresponse.Body)
	if err != nil {
		return err, nil
	}

	// Unmarshal it into our document object
	err = json.Unmarshal(mybody, &docObj)
	if err != nil {
		return err, nil
	}

	return nil, docObj
}

func (box *BoxApi) GetSession(documentId string) (err error, sessObj *SessionObject) {

	encoded, err := json.Marshal(struct {
		DocumentId string `json:"document_id"`
	}{
		documentId,
	})

	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewBuffer(encoded)

	// Create our new request
	mprequest, err := http.NewRequest(
		"POST", box.SessionUrl, body)
	mprequest.Header.Set("Authorization", "Token "+box.ApiKey)
	mprequest.Header.Set("Content-Type", "application/json")
	mprequest.Header.Set("Content-Length", string(body.Len()))

	if err != nil {
		return err, nil
	}

	// Perform the request, grab response and close
	client := &http.Client{}
	mpresponse, err := client.Do(mprequest)
	if err != nil {
		return err, nil
	}
	defer mpresponse.Body.Close()

	// Read the response into a byte slice
	mybody, err := ioutil.ReadAll(mpresponse.Body)
	if err != nil {
		return err, nil
	}

	// Unmarshal it into our document object
	err = json.Unmarshal(mybody, &sessObj)
	if err != nil {
		return err, nil
	}

	return nil, sessObj
}

func (box *BoxApi) GetViewerURL(sessionId string) string {

	return fmt.Sprintf(
		"https://view-api.box.com/1/sessions/%s/view?theme=light",
		sessionId)
}

func NewBoxApi(key string, fileLocation string) *BoxApi {

	return &BoxApi{
		ApiKey:       key,
		FileLocation: fileLocation,
		MultipartUrl: "https://upload.view-api.box.com/1/documents",
		DocumentUrl:  "https://view-api.box.com/1/documents",
		SessionUrl:   "https://view-api.box.com/1/sessions",
	}
}
