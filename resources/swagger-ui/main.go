package main

import (
	"github.com/miiy/goc/resources"
	"io/fs"
	"log"
	"net/http"
)

var swaggerInitializerTpl = `
window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    // url: "https://petstore.swagger.io/v2/swagger.json",
	urls: [{url: "https://petstore.swagger.io/v2/swagger.json", name:"petstore"}, {url: "http://localhost:8101/auth.swagger.json", name:"auth"}],
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

type swaggerInitializer struct {
}

func main() {
	mux := http.NewServeMux()
	subFS, err := fs.Sub(resources.SwaggerUI, "swagger-ui/dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(subFS)))

	s := http.Server{
		Addr:    "0.0.0.0:8101",
		Handler: mux,
	}
	log.Fatal(s.ListenAndServe())
}
