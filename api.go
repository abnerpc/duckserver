package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setAPIRoute(route string) string {
	return fmt.Sprintf("/api/%s/", route)
}

func BuildAPIHandlers() {

	changePassword := http.HandlerFunc(changePasswordHandler)
	addUser := http.HandlerFunc(addUserHandler)
	removeUser := http.HandlerFunc(removeUserHandler)

	http.Handle(setAPIRoute("change_password"), AdminSecureMiddleware(changePassword))
	http.Handle(setAPIRoute("add_user"), AdminSecureMiddleware(addUser))
	http.Handle(setAPIRoute("remove_user"), AdminSecureMiddleware(removeUser))

}

type UserPayload struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type ManipulateUser func(string, string) (string, bool)

func isUserValid(p *UserPayload) bool {
	return p.User != "" && p.Password != ""
}

func parsePayload(w http.ResponseWriter, r *http.Request) (*UserPayload, bool) {
	var p UserPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil || !isUserValid(&p) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return nil, false
	}
	return &p, true
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	manipulateUserHandler(Config.ChangeUserPassword, w, r)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	manipulateUserHandler(Config.AddUser, w, r)
}

func manipulateUserHandler(f ManipulateUser, w http.ResponseWriter, r *http.Request) {
	payload, ok := parsePayload(w, r)
	if !ok {
		return
	}

	msg, ok := f(payload.User, payload.Password)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func removeUserHandler(w http.ResponseWriter, r *http.Request) {}
