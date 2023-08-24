package main

import (
	"flag"
	"github.com/miiy/goc/examples/apiserver/gen"
	"github.com/miiy/goc/openapi"
	"io"
	"io/fs"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", "127.0.0.1:8090", "-addr 127.0.0.1:8090")
)

func main() {
	flag.Parse()

	// scan api files
	urls, err := openapi.ScanApiFiles(gen.OpenAPIFS)
	if err != nil {
		panic(err)
	}
	swaggerInitializerTpl := openapi.ParseTemplate(urls)

	// new multiplexer
	mux := http.NewServeMux()

	// openapi
	mux.Handle("/openapi/", http.FileServer(http.FS(gen.OpenAPIFS)))

	// swagger-ui
	subFS, err := fs.Sub(openapi.SwaggerUIFS, "swagger-ui")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	// custom swagger-initializer.js
	mux.HandleFunc("/swagger-initializer.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		io.WriteString(w, swaggerInitializerTpl)
	})

	// serve http
	s := http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	log.Printf("Serving http://%s\n", *addr)
	log.Fatal(s.ListenAndServe())
}
