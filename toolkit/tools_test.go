package toolkit

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
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
