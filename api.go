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

	listAccessKeys := http.HandlerFunc(listAccessKeysHandler)
	changeAccessKey := http.HandlerFunc(changeAccessKeyHandler)
	addAccessKey := http.HandlerFunc(addAccessKeyHandler)
	deleteAccessKey := http.HandlerFunc(deleteAccessKeyHandler)

	http.Handle(setAPIRoute("list_keys"), AdminSecureMiddleware(listAccessKeys))
	http.Handle(setAPIRoute("change_key"), AdminSecureMiddleware(changeAccessKey))
	http.Handle(setAPIRoute("add_key"), AdminSecureMiddleware(addAccessKey))
	http.Handle(setAPIRoute("delete_key"), AdminSecureMiddleware(deleteAccessKey))

}

func changeAccessKeyHandler(w http.ResponseWriter, r *http.Request) {

	var keys struct {
		OldKey string `json:"old_key"`
		NewKey string `json:"new_key"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&keys)
	if err != nil || keys.OldKey == "" || keys.NewKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}
	msg, ok := Config.ChangeAccessKey(keys.OldKey, keys.NewKey)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)

}

func addAccessKeyHandler(w http.ResponseWriter, r *http.Request) {

	var access struct {
		Key      string `json:"access_key"`
		UserType string `json:"user_type"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&access)
	if err != nil || access.Key == "" || access.UserType == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := Config.AddAccessKey(access.Key, access.UserType)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func deleteAccessKeyHandler(w http.ResponseWriter, r *http.Request) {

	var access struct {
		Key string `json:"access_key"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&access)
	if err != nil || access.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
		return
	}

	msg, ok := Config.DeleteAccessKey(access.Key)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	fmt.Fprintln(w, msg)
}

func listAccessKeysHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := Config.ListAccessKeys()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	result, err := json.Marshal(keys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, string(result))
}
