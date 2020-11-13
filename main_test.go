package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var usersTest = []User{
	User{1, "Jacopo", "jacopo@gmail.com"},
}

func TestUsersIndex(t *testing.T) {
	ts := httptest.NewServer(routeHandler{routes})
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/users", ts.URL))
	if err != nil {
		log.Fatal(err)
	}
	usersJSON, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got %v, want %v", res.StatusCode, http.StatusOK)
	}
	var users []User
	err = json.Unmarshal(usersJSON, &users)
	if err != nil {
		log.Fatal(err)
	}
	for i, user := range users {
		if user != usersTest[i] {
			t.Errorf("got %v, want %v", users, usersTest)
		}
	}
}

func TestUsersIndexInvalidMethod(t *testing.T) {
	ts := httptest.NewServer(routeHandler{routes})
	defer ts.Close()

	res, err := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("got %v, want %v", res.StatusCode, http.StatusNotFound)
	}
}

var user1 = User{1, "Jacopo", "jacopo@gmail.com"}

var userTests = []struct {
	id         uint
	statusCode int
	user       User
}{
	{
		id:         1,
		statusCode: http.StatusOK,
		user:       user1,
	},
	{
		id:         2,
		statusCode: http.StatusNotFound,
		user:       User{},
	},
}

func TestUserGet(t *testing.T) {
	ts := httptest.NewServer(routeHandler{routes})
	defer ts.Close()

	for _, tt := range userTests {
		res, err := http.Get(fmt.Sprintf("%s/users/%d", ts.URL, tt.id))
		if err != nil {
			log.Fatal(err)
		}
		userJSON, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode != tt.statusCode {
			t.Errorf("got %v, want %v", res.StatusCode, tt.statusCode)
		}
		user := User{}
		if len(userJSON) > 0 {
			err = json.Unmarshal(userJSON, &user)
			if err != nil {
				log.Fatal(err)
			}
		}
		if user != tt.user {
			t.Errorf("got %v, want %v", user, tt.user)
		}
	}
}

func TestUserUpdate(t *testing.T) {
	ts := httptest.NewServer(routeHandler{routes})
	defer ts.Close()

	reqBody := map[string]string{
		"name":  "new name",
		"email": "new@email.com",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	resp, err := http.Post(fmt.Sprintf("%s/users/%d", ts.URL, 1), "application/json", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %v, want %v", resp.StatusCode, http.StatusOK)
	}
	userJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	user := User{}
	if len(userJSON) > 0 {
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			log.Fatal(err)
		}
	}
	expectedUser := User{1, "new name", "new@email.com"}
	if user != expectedUser {
		t.Errorf("got %v, want %v", user, expectedUser)
	}
}
