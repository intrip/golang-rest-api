package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func usersIndex(w http.ResponseWriter, req *http.Request) {
	usersJSON, err := json.Marshal(users)
	if err != nil {
		errorMsg(err, w, req)
		return
	}
	fmt.Fprintf(w, string(usersJSON))
}

func usersGet(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)
	user, _ := getUser(id)
	if user == nil {
		notFound(w, req)
		return
	}
	fmt.Fprintf(w, user.RenderJSON())
}

func usersUpdate(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)
	user, _ := getUser(id)
	if user == nil {
		notFound(w, req)
		return
	}
	userJSON, _ := ioutil.ReadAll(req.Body)
	userUpdate := User{}
	err := json.Unmarshal(userJSON, &userUpdate)
	if err != nil {
		errorMsg(err, w, req)
		return
	}
	user.Email = userUpdate.Email
	user.Name = userUpdate.Name
	fmt.Fprintf(w, user.RenderJSON())
}

func usersCreate(w http.ResponseWriter, req *http.Request) {
	userJSON, _ := ioutil.ReadAll(req.Body)
	userCreate := User{}
	err := json.Unmarshal(userJSON, &userCreate)
	if err != nil {
		errorMsg(err, w, req)
		return
	}
	users = append(users, userCreate)
	fmt.Fprintf(w, userCreate.RenderJSON())
}

func usersDelete(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)
	user, idx := getUser(id)
	if user == nil {
		notFound(w, req)
		return
	}
	users[idx] = users[len(users)-1]
	users = users[:len(users)-1]
	w.WriteHeader(http.StatusNoContent)
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// APIError message
type APIError struct {
	message string
}

func errorMsg(errStr error, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	errorJSON, _ := json.Marshal(APIError{errStr.Error()})
	fmt.Fprintf(w, string(errorJSON))
}
