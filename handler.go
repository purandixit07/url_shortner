package main

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

func mapHandler(pathsToURLs map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToURLs[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func yamlHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//1. parse the yaml
	pathURLs, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	//2. convert yaml array into map
	pathsToURLs := buildMap(pathURLs)
	//3. return a maphandler using the map
	return mapHandler(pathsToURLs, fallback), nil
}

func buildMap(pathUrls []pathURL) map[string]string {
	pathsToURLs := make(map[string]string)
	for _, item := range pathUrls {
		pathsToURLs[item.Path] = item.URL
	}
	return pathsToURLs
}

func parseYaml(data []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
