package static

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"fmt"
)

//go:embed "index.js" "style.css"
var StaticFiles embed.FS

type StaticFileHashes map[string]string

var staticFileHashes StaticFileHashes = make(StaticFileHashes)

func GenerateHashesForStaticFiles() {
	dirs, err := StaticFiles.ReadDir(".")

	if err != nil {
		panic(err)
	}

	for _, file := range dirs {
		if file.IsDir() {
			continue
		}

		contents, err := StaticFiles.ReadFile(file.Name())

		if err != nil {
			panic(err)
		}

		h := sha256.New()
		if _, err := h.Write(contents); err != nil {
			panic(err)
		}

		sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
		staticFileHashes[file.Name()] = sha
	}

}

func getHashForFile(filename string) string {
	hash, ok := staticFileHashes[filename]
	if !ok {
		return ""
	}

	return hash
}

func FilePath(filename string) string {
	return fmt.Sprintf("/static/%s?%s", filename, getHashForFile(filename))
}
