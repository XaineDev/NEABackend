package util

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"strconv"
)

func GetImageAsBase64String(imagePath string) (string, error) {
	// open the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}

	// read the image file into a buffer
	imageBuffer := new(bytes.Buffer)
	_, err = io.Copy(imageBuffer, imageFile)
	if err != nil {
		return "", err
	}

	// convert the buffer to a base64 string
	base64String := base64.StdEncoding.EncodeToString(imageBuffer.Bytes())

	return base64String, nil
}

func SaveImageFromBase64String(imageDate string, user int) error {
	// convert the base64 string to a byte array
	imageBytes, err := base64.StdEncoding.DecodeString(imageDate)
	if err != nil {
		return err
	}

	// write the byte array to a file
	err = os.WriteFile(strconv.Itoa(user)+".png", imageBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
