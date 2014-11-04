package boxapi

import "testing"

const key string = "My Super Unique Key"

func TestKeyIsSet(t *testing.T) {

	box := NewBoxApi(key)
	if box.ApiKey != key {
		t.Errorf("Key not set error got %s, want %s", box.ApiKey, key)
	}
}

func TestCanGenerateUniqueFilename(t *testing.T) {

	box := NewBoxApi(key)
	filePath := "http://www.mmta.co.uk/newsletter/crucible"
	expected := "6588e5fba5a3636ba6249eb72246f616"

	if unique := box.generateUniqueFilename(filePath); unique != expected {
		t.Errorf("Filepath %s doesn't match expected: %s", unique, expected)
	}
}
