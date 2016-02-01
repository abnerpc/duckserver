package main

import "net/http"

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = "breakpoint"
		key := r.Header.Get("Authorization")
		if _, ok := Config.AdminKeys[key]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not authorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
