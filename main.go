package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// User is the exposed resource
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RenderJSON render user as JSON
func (u User) RenderJSON() string {
	userJSON, err := json.Marshal(u)
	if err != nil {
		errorJSON, _ := json.Marshal(APIError{err.Error()})
		return string(errorJSON)
	}
	return string(userJSON)
}

// APIError message
type APIError struct {
	message string
}

var users = []User{
	User{1, "Jacopo", "jacopo@gmail.com"},
}

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
	user := getUser(id)
	if user == nil {
		notFound(w, req)
		return
	}
	fmt.Fprintf(w, user.RenderJSON())
}

func usersUpdate(w http.ResponseWriter, req *http.Request) {
	id := extractID(req)
	user := getUser(id)
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

func getUser(id uint) *User {
	for _, u := range users {
		if u.ID == uint(id) {
			return &u
		}
	}
	return nil
}

func notFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func errorMsg(errStr error, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	errorJSON, _ := json.Marshal(APIError{errStr.Error()})
	fmt.Fprintf(w, string(errorJSON))
}

func serveHTTP() {
	mux := http.NewServeMux()
	mux.Handle("/", routeHandler{routes})
	http.ListenAndServe(":8000", mux)
}

func main() {
	serveHTTP()
}

/* TODO
*  - finish Crud for in memory API
*     - CREATE
*     - DELETE
*  - refactor: packages extraction etc
*  - use a DB
 */
