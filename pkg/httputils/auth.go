package httputils

import (
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "net/http"
)

func unauthorized(w http.ResponseWriter) {
    w.Header().Set("WWW-Authenticate", "Basic")
    w.WriteHeader(http.StatusUnauthorized)
}

func NewBasicAuthMiddleware(users map[string]string, restrict bool) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        if !restrict {
            return next
        }

        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestUser, requestPass, _ := r.BasicAuth()
            pass, ok := users[requestUser]
            if !ok {
                unauthorized(w)
                return
            }

            if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(requestPass)); err != nil {
                unauthorized(w)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
