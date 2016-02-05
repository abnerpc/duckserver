package main

import "net/http"

func SecureMiddleware(next http.Handler, userTypes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		userType, ok := Config.AccessKeys[key]
		if ok {
			ok = false
			for _, ut := range userTypes {
				if userType == ut {
					ok = true
					break
				}
			}
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not authorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminSecureMiddleware(next http.Handler) http.Handler {
	return SecureMiddleware(next, Admin)
}

func UserSecureMiddleware(next http.Handler) http.Handler {
	return SecureMiddleware(next, Admin, User)
}
