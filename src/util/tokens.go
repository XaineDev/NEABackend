package util

import "github.com/google/uuid"

func CreateToken() (string, error) {
	id, err := uuid.NewUUID()
	return id.String(), err
}
