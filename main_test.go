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

var ts *httptest.Server

func setup() {
	ts = httptest.NewServer(routeHandler{routes})
	users = usersTest
}

func teardown() {
	ts.Close()
}

func TestUsersIndex(t *testing.T) {
	setup()
	defer teardown()

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
	setup()
	defer teardown()

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users", ts.URL), nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("got %v, want %v", res.StatusCode, http.StatusNotFound)
	}
}

var userTests = []struct {
	id         uint
	statusCode int
	user       User
}{
	{
		id:         1,
		statusCode: http.StatusOK,
		user:       User{1, "Jacopo", "jacopo@gmail.com"},
	},
	{
		id:         2,
		statusCode: http.StatusNotFound,
		user:       User{},
	},
}

func TestUsersGet(t *testing.T) {
	setup()
	defer teardown()

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

func TestUsersCreate(t *testing.T) {
	setup()
	defer teardown()

	createUser := User{2, "new name", "new@email.com"}
	createUserJSON, _ := json.Marshal(createUser)

	resp, err := http.Post(fmt.Sprintf("%s/users", ts.URL), "application/json", bytes.NewBuffer(createUserJSON))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %v, want %v", resp.StatusCode, http.StatusOK)
	}
	userJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	user := User{}
	if len(userJSON) > 0 {
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			log.Fatal(err)
		}
	}
	if user != createUser {
		t.Errorf("got %v, want %v", user, createUser)
	}
	if len(users) != 2 {
		t.Errorf("users length is wrong: got %d, want %d", len(users), 2)
	}
}

func TestUsersUpdate(t *testing.T) {
	setup()
	defer teardown()

	updateUser := User{1, "new name", "new@email.com"}
	updateUserJSON, _ := json.Marshal(updateUser)

	resp, err := http.Post(fmt.Sprintf("%s/users/%d", ts.URL, updateUser.ID), "application/json", bytes.NewBuffer(updateUserJSON))
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
	if user != updateUser {
		t.Errorf("got %v, want %v", user, updateUser)
	}
}
