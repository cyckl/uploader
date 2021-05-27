// Auth / security module I guess
// At least I use bcrypt lmao

package main

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// Don't know what this really means but it's hashing related
var cost = 14

type AuthFile struct {
	User	string
	Secret	string
}

// Authenticate!
func auth(w http.ResponseWriter, r *http.Request) error {
	// Read auth file
	data, err := ioutil.ReadFile("./auth.json")
	if err != nil {
		return errors.New("failed to open auth file")
	}
	
	// Link JSON data slice to auth struct
	var auth AuthFile
	err = json.Unmarshal(data, &auth)
	if err != nil {
		return errors.New(fmt.Sprintf("could not unmarshal JSON data: %v\n", err))
	}

	u, s, ok := r.BasicAuth()
	// Check if auth in header was formed correctly
	if !ok {
		w.WriteHeader(401)
		return errors.New("malformed")
	}
	
	// Check credentials
	if u != auth.User {
		w.WriteHeader(401)
		return errors.New("incorrect credentials")
	}
	if pwdCheck(s, auth.Secret) != true {
		w.WriteHeader(401)
		return errors.New("incorrect credentials")
	}
	
	return nil
}

// Set a new user in the config
func setUser(u string) {
	// Read auth file
	data, err := ioutil.ReadFile("./auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to open file: %v\n", err)
		return
	}
	
	// Link JSON data slice to Auth struct
	var auth AuthFile
	err = json.Unmarshal(data, &auth)
	if err != nil {
		log.Printf("[Auth] Could not unmarshal JSON data: %v\n", err)
		return
	}
	
	// Set new user in field
	auth.User = string(u)
	
	// Re-encode in JSON
	authNew, err := json.Marshal(auth)
	if err != nil {
		log.Printf("[Auth] Could not marshal JSON data: %v\n", err)
		return
	}
	
	// Save to file
	err = ioutil.WriteFile("./auth.json", authNew, 0644)
	if err != nil {
		log.Printf("[Auth] Could not write to file: %v\n", err)
		return
	}
}

// Set a new secret in the config
func setSecret(s string) {
	// Hash the new secret
	sBytes, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		log.Printf("[Auth] Could not hash new secret: %v\n", err)
		return
	}
	
	// Read auth file
	data, err := ioutil.ReadFile("auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to open file: %v\n", err)
		return
	}
	
	// Link JSON data slice to Auth struct
	var auth AuthFile
	err = json.Unmarshal(data, &auth)
	if err != nil {
		log.Printf("[Auth] Could not unmarshal JSON data: %v\n", err)
		return
	}
	
	// Set new secret in field
	auth.Secret = string(sBytes)
	
	// Re-encode in JSON
	authNew, err := json.Marshal(auth)
	if err != nil {
		log.Printf("[Auth] Could not marshal JSON data: %v\n", err)
		return
	}
	
	// Save to file
	err = ioutil.WriteFile("./auth.json", authNew, 0644)
	if err != nil {
		log.Printf("[Auth] Could not write to file: %v\n", err)
		return
	}
}

// This is probably a pretty important bit
func pwdCheck(p, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
