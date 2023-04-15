package util

import (
	"encoding/json"
	"net/http"
)

func ToJson(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func RespondWithJson(writer http.ResponseWriter, v interface{}) error {
	data, err := ToJson(v)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ParseJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
