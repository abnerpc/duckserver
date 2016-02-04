package main

import "net/http"

func SecureMiddleware(allowedKeys map[string]string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if _, ok := allowedKeys[key]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not authorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminSecureMiddleware(next http.Handler) http.Handler {
	return SecureMiddleware(Config.AdminKeys, next)
}

func UserSecureMiddleware(next http.Handler) http.Handler {
	if Config.UserKeys == nil {
		return AdminSecureMiddleware(next)
	}
	keys := make(map[string]string)
	for k, v := range Config.AdminKeys {
		keys[k] = v
	}
	for k, v := range Config.UserKeys {
		keys[k] = v
	}
	return SecureMiddleware(keys, next)
}
