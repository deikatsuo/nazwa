package misc

import (
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

// Base64ToFileWithData mengembalikan kembali string base64 kedalam bentuk file
// dir: lokasi file
// b64: string base64
// info: informasi file
func Base64ToFileWithData(dir string, b64 string, info string) (string, error) {
	ext := Base64GetExtensionFromData(info)
	fileName := uuid.New().String() + "." + ext
	decode, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println("ERROR: image.go Base64ToFileWithData() Meng decode string base64")
		return "", err
	}

	fs, err := os.Create(dir + fileName)
	if err != nil {
		log.Println("ERROR: image.go Base64ToFileWithData() Membuat file")
		return "", err
	}

	defer fs.Close()

	if _, err := fs.Write(decode); err != nil {
		log.Println("ERROR: image.go Base64ToFileWithData() Menulis kedalam file")
		return "", err
	}

	if err := fs.Sync(); err != nil {
		log.Println("ERROR: image.go Base64ToFileWithData() Sinkron file")
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
