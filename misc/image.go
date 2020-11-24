package misc

import (
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

// Base64ToImageFileWithData mengembalikan kembali string base64 kedalam bentuk file
func Base64ToImageFileWithData(b64 string, data string) (string, error) {
	ext := Base64GetExtensionFromData(data)
	fileName := uuid.New().String() + "." + ext
	decode, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println("ERROR: image.go Base64ToImageFileWithData() Meng decode string base64")
		return "", err
	}

	fs, err := os.Create("./upload/profile/" + fileName)
	if err != nil {
		log.Println("ERROR: image.go Base64ToImageFileWithData() Membuat file")
		return "", err
	}

	defer fs.Close()

	if _, err := fs.Write(decode); err != nil {
		log.Println("ERROR: image.go Base64ToImageFileWithData() Menulis kedalam file")
		return "", err
	}

	if err := fs.Sync(); err != nil {
		log.Println("ERROR: image.go Base64ToImageFileWithData() Sinkron file")
		return "", err
	}

	return fileName, nil
}

// Base64GetExtensionFromData mengembalikan extensi dari data
func Base64GetExtensionFromData(data string) string {
	gext := strings.ReplaceAll(data, "data:image/", "")
	gext = strings.ReplaceAll(gext, ";base64", "")
	return gext
}
