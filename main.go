//
//	I can't tell if this is good code or bad code...
//	That probably means that it's bad code
//

package main

import (
	"errors"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"
)
var (
	// Set up flags and their defaults
	size = flag.Int64("m", 10, "The max file size in MB")
	port = flag.String("p", "8080", "The port to bind to")
	dir = flag.String("d", "", "Location to save files in")
	host = flag.String("w", "", "Public-facing URL for server")
	
	// Handle setting new creds but default to empty
	newUser = flag.String("u", "", "Set a new auth username")
	newSecret = flag.String("s", "", "Set a new auth secret")
)

func main() {
	// Pass in flags
	flag.Parse()
	
	// Check if there are new credentials being set
	if *newUser != "" {
		setUser(*newUser)
		log.Printf("Setting new username: %v\n", *newUser)
	}
	if *newSecret != "" {
		setSecret(*newSecret)
		log.Printf("Setting new user secret\n")
	}
	
	// Delegate URL endpoint and call function
	http.HandleFunc("/upload", upload)
	
	// Bind to port
	log.Printf("Listening on port %v\n", *port)
	err := http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

// Upload to server
func upload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting new upload from %v\n", r.RemoteAddr)
	
	// Check authentication
	err := auth(w, r)
	if err != nil {
		log.Printf("[Error] Authentication failed: %v\n", err)
		return
	}
	
	// Parse form with max file size in MB
	err = r.ParseMultipartForm(*size << 20)
	if err != nil {
		log.Printf("[Error] Failed to parse multipart form: %v\n", err)
		return
	}
	
	// Return file data for the HTML tag "data"
	file, header, err := r.FormFile("data")
	if err != nil {
		log.Printf("[Error] Failed to get file from uploader: %v\n", err)
		return
	}
	defer file.Close()
	
	// Read upload to bytestream
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[Error] Failed to get raw byte data from upload: %v\n", err)
		return
	}
	
	// Save to file
	name, err := save(data, header)
	if err != nil {
		log.Printf("[Error] Failed to save byte data to file: %v\n", err)
		return
	}
	
	// Log successful upload
	log.Printf("Saved %v (%v bytes) from %v\n", name, header.Size, r.RemoteAddr)
	
	// Send back response with URL + file name
	fmt.Fprintln(w, *host + name)
}

// Yes
func save(data []byte, header *multipart.FileHeader) (string, error) {
	// Generate name
	name, err := nameGen(header.Filename)
	if err != nil {
		return "", errors.New(fmt.Sprintf("random name generation failed: %v\n", err))
	}
	
	// Set save location
	loc := *dir + name
	
	// Save that bytestream to a file with 644 perms
	err = ioutil.WriteFile(loc, data, 0644)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to save raw byte data as file: %v\n", err))
	}
	
	return name, nil
}

func nameGen(file string) (string, error) {
	// Read word file
	data, err := ioutil.ReadFile("./words.json")
	if err != nil {
		return "", errors.New("failed to open word file")
	}
	
	// Link JSON data slice to word slice
	var words []string
	err = json.Unmarshal(data, &words)
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not unmarshal JSON data: %v\n", err))
	}

	// Shuffle words in word list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	// Get first three entries of shuffled array
	gen := words[0] + words[1] + words[2]
	
	// Filename "edge cases" (they're common but just *special*)
	var name string
	if strings.Contains(file, strings.ToLower(".tar.")) {
		name = string(gen) + ".tar" + path.Ext(file)
	} else if path.Ext(file) == "" {
		name = string(gen) + ".bin"
	} else {
		name = string(gen) + path.Ext(file)
	}
	
	return name, nil
}
