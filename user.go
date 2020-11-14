package main

import (
	"encoding/json"
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

func getUser(id uint) *User {
	for _, u := range users {
		if u.ID == uint(id) {
			return &u
		}
	}
	return nil
}
