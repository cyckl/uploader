//
//	I can't tell if this is good code or bad code...
//	That probably means that it's bad code
//

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	}
	if *newSecret != "" {
		setSecret(*newSecret)
	}
	
	// Delegate URL endpoint and call function
	http.HandleFunc("/upload", upload)
	
	// Bind to port
	log.Printf("[Status] Attempting bind to port %v\n", *port)
	err := http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

// Upload and save file
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
	file, handler, err := r.FormFile("data")
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
	
	// Set save location
	name := nameGen(handler.Filename, 5)
	loc := *dir + name
	
	// Save that bytestream to a file with 644 perms
	err = ioutil.WriteFile(loc, data, 0644)
	if err != nil {
		log.Printf("[Error] Failed to save raw byte data as file: %v\n", err)
		return
	}
	
	// Log successful upload
	log.Printf("Saved %v (%v bytes) from %v\n", name, handler.Size, r.RemoteAddr)
	
	// Send back response with URL + file name
	fmt.Fprintln(w, *host + name)
}

func nameGen(orig string, l int) string {
	// Generate random file name
	rand.Seed(time.Now().UnixNano())
	var char = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	gen := make([]rune, l)
	for i := range gen {
		gen[i] = char[rand.Intn(len(char))]
	}
	
	// Filename "edge cases" (they're common but just *special*)
	var name string
	if strings.Contains(orig, strings.ToLower(".tar.")) {
		name = string(gen) + ".tar" + path.Ext(orig)
	} else if path.Ext(orig) == "" {
		name = string(gen) + ".bin"
	} else {
		name = string(gen) + path.Ext(orig)
	}
	
	// Return the new path
	return name
}
