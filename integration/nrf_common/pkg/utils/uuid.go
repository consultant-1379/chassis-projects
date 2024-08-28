package utils

import (
	"fmt"

	gouuid "github.com/satori/go.uuid"
)

func GetUUIDString() (string, error) {
	uuid, err := gouuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", uuid), err
}
