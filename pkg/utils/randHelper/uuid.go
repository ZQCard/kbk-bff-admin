package randHelper

import (
	"fmt"

	uuidpkg "github.com/google/uuid"
)

func GenUUID() (string, error) {
	id, err := uuidpkg.NewRandom()
	if err != nil {
		return "", fmt.Errorf("生成 uuid 失败: %w", err)
	}
	return id.String(), nil
}

func GetUUID() string {
	s, err := GenUUID()
	if err != nil {
		return ""
	}
	return s
}
