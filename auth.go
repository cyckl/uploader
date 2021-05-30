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

func auth(w http.ResponseWriter, r *http.Request) error {
	// Read in data from auth file
	a, err := readAuth("./auth.json")
	if err != nil {
		return errors.New(fmt.Sprintf("could not read from auth file: %v\n", err))
	}

	u, s, ok := r.BasicAuth()
	// Check if auth in header was formed correctly
	if !ok {
		w.WriteHeader(401)
		return errors.New("malformed")
	}
	
	// Check credentials
	if u != a.User {
		w.WriteHeader(401)
		return errors.New("incorrect credentials")
	}
	if pwdCheck(s, a.Secret) != true {
		w.WriteHeader(401)
		return errors.New("incorrect credentials")
	}
	
	return nil
}

// Set a new user in the config
func setUser(u string) {
	// Read in data from auth file
	a, err := readAuth("./auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to read in auth data: %v\n", err)
		return
	}
	
	// Set new user in field
	a.User = string(u)
	
	// save data to auth file
	err = saveAuth(a, "./auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to write auth data: %v\n", err)
		return
	}
}

// Set a new secret in the config
func setSecret(s string) {
	// Hash the new secret
	bs, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		log.Printf("[Auth] Could not hash new secret: %v\n", err)
		return
	}
	
	// Read in data from auth file
	a, err := readAuth("./auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to read in auth data: %v\n", err)
		return
	}
	
	// Set new secret in field
	a.Secret = string(bs)
	
	// save data to auth file
	err = saveAuth(a, "./auth.json")
	if err != nil {
		log.Printf("[Auth] Failed to write auth data: %v\n", err)
		return
	}
}

func pwdCheck(p, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}

func readAuth(path string) (auth AuthFile, err error) {
	// Read auth file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return auth, errors.New("failed to open file")
	}
	
	// Link JSON data slice to Auth struct
	err = json.Unmarshal(data, &auth)
	if err != nil {
		return auth, errors.New(fmt.Sprintf("could not unmarshal JSON data: %v\n", err))
	}
	
	return auth, nil
}

func saveAuth(auth AuthFile, path string) error {
	// Re-encode in JSON
	json, err := json.MarshalIndent(auth, "", "\t")
	if err != nil {
		return errors.New(fmt.Sprintf("could not re-marshal JSON data: %v\n", err))
	}
	
	// Save to file
	err = ioutil.WriteFile(path, json, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("could not write to file: %v\n", err))
	}
	
	return nil
}
