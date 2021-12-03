package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func FromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	var sb strings.Builder
	for _, v := range sum {
		sb.WriteString(fmt.Sprintf("%02x", v))
	}
	return sb.String(), nil
}
