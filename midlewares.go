package main

import (
	"net/http"
	"strings"
)

// GetAuthorizarion split the header Authorization string and returns O_WRONLY
// the token after the prefix
func GetAuthorizarion(auth string) (string, bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return "", false
	}
	return auth[len(prefix):], true
}

// SecureMiddleware verifys the Authorization token send with the request
func SecureMiddleware(next http.Handler, userTypes ...byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorized := false
		key, ok := GetAuthorizarion(r.Header.Get("Authorization"))
		if ok {
			userType, ok := CurrentConfig.AccessKeys[key]
			if ok {
				for _, ut := range userTypes {
					if userType == ut {
						authorized = true
						break
					}
				}
			}
		}
		if !authorized {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"user\"")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not authorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AdminSecureMiddleware verifys the token checking if the user is Admin
func AdminSecureMiddleware(next http.Handler) http.Handler {
	return SecureMiddleware(next, Admin)
}

// UserSecureMiddleware verifys the token checking if the user is Admin or User
func UserSecureMiddleware(next http.Handler) http.Handler {
	return SecureMiddleware(next, Admin, User)
}
