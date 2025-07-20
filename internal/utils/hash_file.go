package utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash/v2"
)

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	hasher := xxhash.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	hashValue := hasher.Sum64()

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, hashValue)

	b64 := base64.StdEncoding.EncodeToString(buf)

	return b64, nil
}
