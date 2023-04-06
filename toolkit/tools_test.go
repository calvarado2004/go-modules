package toolkit

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

// TestTools_RandomString tests the RandomString function
func TestTools_RandomString(t *testing.T) {
	var testTools Tools
	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("Expected 10, got %d", len(s))
	}

}

// TestTools_UploadFiles tests the UploadFiles function
func TestTools_UploadFiles(t *testing.T) {
	var uploadTests = []struct {
		name          string
		allowedTypes  []string
		renameFile    bool
		errorExpected bool
	}{
		{"allowed no rename", []string{"image/jpeg", "image/png"}, false, false},
		{"allowed rename", []string{"image/jpeg", "image/png"}, true, false},
		{"not allowed", []string{"image/png"}, false, true},
		{"not allowed rename", []string{"image/png"}, true, true},
	}

	for _, tt := range uploadTests {
		// setup a pipe to avoid buffering
		pr, pw := io.Pipe()
		// create a new multipart writer
		mw := multipart.NewWriter(pw)
		// create a wait group
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer mw.Close()
			defer wg.Done()

			// create a new form data filed for the file
			part, err := mw.CreateFormFile("file", "./testdata/nyc.jpeg")
			if err != nil {
				t.Errorf("Error creating form file: %v", err)
			}

			file, err := os.Open("./testdata/nyc.jpeg")
			if err != nil {
				t.Errorf("Error opening file: %v", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				t.Errorf("Error decoding image: %v", err)
			}

			err = jpeg.Encode(part, img, nil)
			if err != nil {
				t.Errorf("Error encoding image: %v", err)
			}

		}()

		// read from the pipe which receives the multipart data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", mw.FormDataContentType())

		var testTools Tools

		testTools.AllowedFileTypes = tt.allowedTypes

		// test UploadFiles function
		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", tt.renameFile)
		if err != nil && !tt.errorExpected {
			t.Errorf("Error uploading file: %v", err)
		}

		if !tt.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("File does not exist on the destination directory: %s on test %s", err.Error(), tt.name)
			}

			// remove the file
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if tt.errorExpected && err == nil {
			t.Errorf("%s: Error was expected but none received", tt.name)
		}

		wg.Wait()

	}
}

// TestTools_UploadOneFile tests the UploadOneFile function
func TestTools_UploadOneFile(t *testing.T) {
	var uploadTests = []struct {
		name          string
		allowedTypes  []string
		renameFile    bool
		errorExpected bool
	}{
		{"allowed no rename", []string{"image/jpeg", "image/png"}, false, false},
		{"allowed rename", []string{"image/jpeg", "image/png"}, true, false},
		{"not allowed", []string{"image/png"}, false, true},
		{"not allowed rename", []string{"image/png"}, true, true},
	}

	for _, tt := range uploadTests {
		// setup a pipe to avoid buffering
		pr, pw := io.Pipe()
		// create a new multipart writer
		mw := multipart.NewWriter(pw)

		go func() {
			defer mw.Close()

			// create a new form data filed for the file
			part, err := mw.CreateFormFile("file", "./testdata/nyc.jpeg")
			if err != nil {
				t.Errorf("Error creating form file: %v", err)
			}

			file, err := os.Open("./testdata/nyc.jpeg")
			if err != nil {
				t.Errorf("Error opening file: %v", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				t.Errorf("Error decoding image: %v", err)
			}

			err = jpeg.Encode(part, img, nil)
			if err != nil {
				t.Errorf("Error encoding image: %v", err)
			}

		}()

		// read from the pipe which receives the multipart data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", mw.FormDataContentType())

		var testTools Tools

		testTools.AllowedFileTypes = tt.allowedTypes

		// test UploadOneFile function
		uploadedFiles, err := testTools.UploadOneFile(request, "./testdata/uploads/", tt.renameFile)
		if err != nil && !tt.errorExpected {
			t.Errorf("Error uploading file: %v", err)
		}

		if !tt.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName)); os.IsNotExist(err) {
				t.Errorf("File does not exist on the destination directory: %s on test %s", err.Error(), tt.name)
			}

			// remove the file
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", uploadedFiles.NewFileName))
		}

		if tt.errorExpected && err == nil {
			t.Errorf("%s: Error was expected but none received", tt.name)
		}

	}
}

// TestTools_CreateDirIfNotExist tests the CreateDirIfNotExist function
func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools Tools

	err := testTools.CreateDirIfNotExist("./testdata/mydir")
	if err != nil {
		t.Errorf("Error creating directory: %v", err)
	}

	// remove the directory
	_ = os.Remove("./testdata/mydir")

}

// TestTools_Slugify tests the Slugify function
func TestTools_Slugify(t *testing.T) {
	var slugTests = []struct {
		name          string
		input         string
		expected      string
		errorExpected bool
	}{
		{"valid string", "This is a test", "this-is-a-test", false},
		{"valid string with numbers", "This is a test 123", "this-is-a-test-123", false},
		{"valid string with special characters", "This is a test 123 !@#$%^&*()_+{}|:<>?[]\\;',./", "this-is-a-test-123", false},
		{"empty string", " ", "this-is-a-test-123", true},
		{"japanese strings", "こんにちは世界", "kon-nichi-ha-se-kai", true},
	}

	for _, tt := range slugTests {
		var testTools Tools

		slug, err := testTools.Slugify(tt.input)
		if slug != tt.expected && !tt.errorExpected {
			t.Errorf("%s: Expected %s but got %s with error %s", tt.name, tt.expected, slug, err)
		}

		if !tt.errorExpected && slug != tt.expected {
			t.Errorf("%s: Expected %s but got %s", tt.name, tt.expected, slug)
		}

		if tt.errorExpected && err == nil {
			t.Errorf("%s: Error was expected but none received", tt.name)
		}
	}
}

// TestTools_DownloadStaticFile tests the DownloadStaticFile function
func TestTools_DownloadStaticFile(t *testing.T) {
	rr := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)

	var testTools Tools

	testTools.DownloadStaticFile(rr, request, "./testdata", "nyc.jpeg", "nyc.jpeg")

	res := rr.Result()

	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "226989" {
		t.Errorf("Wrong content length, expected %s", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Type"][0] != "image/jpeg" {
		t.Errorf("Wrong content type, expected %s", res.Header["Content-Type"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"nyc.jpeg\"" {
		t.Errorf("Wrong content disposition, expected %s", res.Header["Content-Disposition"][0])
	}

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

}

// TestTools_ReadJSON tests the ReadJSON function
func TestTools_ReadJSON(t *testing.T) {
	var jsonTests = []struct {
		name          string
		json          string
		errorExpected bool
		maxSize       int
		allowUnknown  bool
	}{
		{"valid json", `{"name": "John", "age": 30}`, false, 4096, false},
		{"invalid json", `{"name": "John", "age": 30`, true, 4096, false},
		{"json with unknown fields", `{"name": "John", "age": 30, "city": "New York"}`, true, 4096, false},
		{"json with unknown fields allowed", `{"name": "John", "age": 30, "city": "New York"}`, false, 4096, true},
		{"json with max size", `{"name": "John", "age": 30}`, true, 10, false},
		{"json with incorrect type", `{"name": "John", "age": "30"}`, true, 4096, false},
		{"json with two json objects", `{"name": "John", "age": 30}{"name": "John", "age": 30}`, true, 4096, false},
		{"empty json", ``, true, 4096, false},
		{"json with missing fieldname", `{:jack, "age": 30}`, true, 4096, false},
		{"json with missing fieldvalue", `{"name": "John", "age": }`, true, 4096, false},
		{"not json", `not json`, true, 4096, false},
	}

	var testTools Tools

	for _, tt := range jsonTests {

		testTools.MaxJSONSize = tt.maxSize
		testTools.AllowUnknownFields = tt.allowUnknown

		var decodedJSON struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(tt.json)))
		if err != nil {
			t.Log("Error:", err)
		}

		rr := httptest.NewRecorder()

		err = testTools.ReadJSON(rr, req, &decodedJSON)

		if tt.errorExpected && err == nil {
			t.Errorf("%s: Error was expected but none received", tt.name)
		}

		if !tt.errorExpected && err != nil {
			t.Errorf("%s: Error was not expected but received: %v", tt.name, err.Error())
		}

		req.Body.Close()

	}
}

// TestTools_WriteJSON tests the WriteJSON function
func TestTools_WriteJSON(t *testing.T) {
	var testTools Tools

	rr := httptest.NewRecorder()
	payload := JSONResponse{
		Error: false,
		Msg:   "test",
	}

	headers := make(http.Header)

	headers.Add("Content-Type", "application/json")

	err := testTools.WriteJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("Error writing json: %v", err)
	}

}
