//go:build !dev
// +build !dev

package static

import (
	"fmt"
	"net/http"
)

func StaticFileServer() http.Handler {
	fmt.Println("Embedding in prod!")
	return http.FileServerFS(StaticFiles)

}
