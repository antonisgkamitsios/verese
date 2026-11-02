//+build dev
//go:build dev
// +build dev

package static

import (
	"fmt"
	"net/http"
)

func StaticFileServer() http.Handler {
	fmt.Println("building static files for development")
	return http.FileServer(http.Dir("static"))
}
