package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jatin-malik/url-shortener/urlshort"
)

// Command line flags
var yaml_file = flag.String("yaml", "", "yaml file for path and url mappings")
var json_file = flag.String("json", "", "json file for path and url mappings")

func ReadFile(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error opening file")
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("error reading file")
	}
	return data, nil
}

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	var handler http.Handler = mapHandler
	if *yaml_file != "" {

		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yaml_data, err := ReadFile(*yaml_file)
		if err != nil {
			log.Fatal(err)
		}
		yamlHandler, err := urlshort.YAMLHandler(yaml_data, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = yamlHandler
	}

	if *json_file != "" {

		// Build the JsonHandler using the mapHandler as the fallback

		json_data, err := ReadFile(*json_file)
		if err != nil {
			log.Fatal(err)
		}
		jsonHandler, err := urlshort.JSONHandler(json_data, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = jsonHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
