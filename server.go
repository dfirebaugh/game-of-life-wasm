package main

import (
	"log"
	"net/http"
	"strings"
)

const dir = "./client"
const wasmdir = "./wasm/bin"

func main() {
	fs := http.FileServer(http.Dir(dir))
	wasmfs := http.FileServer(http.Dir(wasmdir))
	log.Print("Serving " + dir + " on http://localhost:8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		// disable cache control for developement only
		resp.Header().Add("Cache-Control", "no-cache")
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
			wasmfs.ServeHTTP(resp, req)
		}
		fs.ServeHTTP(resp, req)
	}))
}
