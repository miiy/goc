package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/miiy/goc/examples/apiserver/api"
	"github.com/miiy/goc/openapi"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"text/template"
)

type swaggerInitializerUrl struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

var swaggerInitializerTpl = `
window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    // url: "https://petstore.swagger.io/v2/swagger.json",
    urls: {{.Urls}},
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
`

var (
	addr = flag.String("addr", "127.0.0.1:8090", "-addr 127.0.0.1:8090")
)

func main() {
	flag.Parse()

	// scan api docs
	urls, err := scanApiFiles(api.OpenAPIFS)
	if err != nil {
		panic(err)
	}
	swaggerInitializerTpl = parseTemplate(urls)

	mux := http.NewServeMux()
	// docs
	mux.Handle("/docs/", http.FileServer(http.FS(api.OpenAPIFS)))

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

// scanApiFiles scan openapi/*.json file
func scanApiFiles(fsys fs.FS) ([]swaggerInitializerUrl, error) {
	de, err := fs.ReadDir(fsys, "docs")
	if err != nil {
		return nil, err
	}

	var urls []swaggerInitializerUrl
	for _, f := range de {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".json" {
			continue
		}
		urls = append(urls, swaggerInitializerUrl{
			Url:  "/docs/" + f.Name(),
			Name: f.Name(),
		})
	}
	return urls, nil
}

// parseTemplate parse swagger initializer file
func parseTemplate(urls []swaggerInitializerUrl) string {
	urlsJson, err := json.Marshal(urls)
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("swaggerInitializerTpl").Parse(swaggerInitializerTpl))
	b := &bytes.Buffer{}
	err = tpl.Execute(b, map[string]interface{}{"Urls": string(urlsJson)})
	if err != nil {
		panic(err)
	}
	return b.String()
}
