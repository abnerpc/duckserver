package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setAPIRoute(route string) string {
	return fmt.Sprintf("/api/%s/", route)
}

// BuildAPIHandlers set up the api handlers
func BuildAPIHandlers() {

	changePassword := http.HandlerFunc(changePasswordKeyHandler)
	addUser := http.HandlerFunc(addUserHandler)
	deleteUser := http.HandlerFunc(deleteUserHandler)

	http.Handle(setAPIRoute("change_password"), AdminSecureMiddleware(changePassword))
	http.Handle(setAPIRoute("add_user"), AdminSecureMiddleware(addUser))
	http.Handle(setAPIRoute("delete_user"), AdminSecureMiddleware(deleteUser))

}

func changePasswordKeyHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		AccessKey   string `json:"access_key"`
		NewPassword string `json:"new_password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.AccessKey == "" || user.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}
	msg, ok := Config.changePassword(user.AccessKey, user.NewPassword)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		AccessKey string `json:"access_key"`
		UserType  byte   `json:"user_type"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.AccessKey == "" || (user.UserType != Admin && user.UserType != User) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := Config.addUser(user.AccessKey, user.UserType)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	var user struct {
		AccessKey string `json:"access_key"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.AccessKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := Config.deleteUser(user.AccessKey)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}
