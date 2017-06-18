package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"os"
)

const (
	apiLinkLibrary	string = "https://registry.hub.docker.com/v2/repositories/library/%s/tags/"
	apiLinkRepo	string = "https://registry.hub.docker.com/v2/repositories/%s/tags/"
)

//--------------------------------------------------------------
type Config struct {
	Type		string		`json:"type"`
	Name		string		`json:"name"`
	Id		int64		`json:"id"`
	Last_updated	string		`json:"last_updated"`
}

//--------------------------------------------------------------
type TagResult struct {
	Name		string 		`json:"name"`
	Id		int64		`json:"id"`
	Last_updated	string		`json:"last_updated"`
}

type Tags struct {
	Detail		string		`json:"detail"`
	Count		int64		`json:"count"`
	Next		string		`json:"next"`
	Previous	string		`json:"previous"`
	Results		[]TagResult	`json:"results"`
}

//--------------------------------------------------------------
func getTags(repo string) Tags {
	var ts Tags
	link := fmt.Sprintf(apiLinkRepo, repo)
	
	// Get data using api registry.hub.docker/v2
	resp, err := http.Get(link)
	if err != nil {
		ts.Detail = err.Error()
		return ts
	}
	
	// Close connection
	defer resp.Body.Close()
	
	// Read data
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ts.Detail = err.Error()
		return ts
	}

	// Unmarshal json
	json.Unmarshal(raw, &ts)
	return ts
}

//--------------------------------------------------------------
func getConfig(fn string) []Config {
	var cf []Config
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	json.Unmarshal(raw, &cf)
	return cf
}

func fmJsonConfig(str string) string {
	str = strings.Replace(str, "\"type\":", "\n\t\"type\":", -1)
	str = strings.Replace(str, "\"name\":", "\n\t\"name\":", -1)
	str = strings.Replace(str, "\"id\":", "\n\t\"id\":", -1)
	str = strings.Replace(str, "\"last_updated\":", "\n\t\"last_updated\":", -1)
	str = strings.Replace(str, "\"}", "\"\n}\n", -1)
	return str
}

func saveConfig(fn string, cf []Config) error {
	// Open file
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	
	// Convert to Json
	bytes, err := json.Marshal(cf)
	if err != nil {
		f.Close()
		return err
	}

	// Write file
	if _, err = f.WriteString(fmJsonConfig(string(bytes))); err != nil {
		f.Close()
		return err
	}
	f. Close()
	return nil
}
