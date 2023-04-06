package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const randonStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Tools is a toolkit for general purpose
type Tools struct {
	MaxFileSize      int64
	AllowedFileTypes []string
}

// RandomString returns a random string of length n
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randonStringSource)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))

		s[i] = r[x%y]
	}

	return string(s)
}

// UploadedFile is a struct to hold uploaded file information
type UploadedFile struct {
	NewFileName  string
	OriginalFile string
	FileSize     int64
}

// UploadOneFile uploads one file to the server
func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {

	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// UploadFiles uploads files to the server and verifies that the file type is allowed
func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := t.CreateDirIfNotExist(uploadDir)
	if err != nil {
		return nil, err
	}

	err = r.ParseMultipartForm(t.MaxFileSize)
	if err != nil {
		return nil, errors.New("the uploaded file is too large")
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {

				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}

				defer infile.Close()

				buff := make([]byte, 512)

				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				allowed := false
				fileType := http.DetectContentType(buff)

				if len(t.AllowedFileTypes) > 0 {
					for _, allowedType := range t.AllowedFileTypes {
						if strings.EqualFold(fileType, allowedType) {
							allowed = true
							break
						}
					}

				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("the file type is not allowed")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				uploadedFile.OriginalFile = hdr.Filename

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				}

				fileSize, err := io.Copy(outfile, infile)
				if err != nil {
					return nil, err
				}

				uploadedFile.FileSize = fileSize

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}

// CreateDirIfNotExist creates a directory if it does not exist
func (t *Tools) CreateDirIfNotExist(path string) error {

	const mode = os.ModeDir | 0755

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

// Slugify returns a slugified string
func (t *Tools) Slugify(s string) (string, error) {

	if s == "" {
		return "", errors.New("string is empty")
	}

	var re = regexp.MustCompile("[^a-zA-Z0-9]+")

	slug := strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")

	if len(slug) == 0 {
		return "", errors.New("slug is empty")
	}

	return slug, nil

}

// DownloadStaticFile downloads a static file
func (t *Tools) DownloadStaticFile(w http.ResponseWriter, r *http.Request, osPath, file, displayName string) {

	filePath := path.Join(osPath, file)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", displayName))

	http.ServeFile(w, r, filePath)

}