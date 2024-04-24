package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (app *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		_, err := app.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		/*out, err := json.Marshal(payload)
		if err != nil {
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}
		r.Header.Set(authorizationPayloadKey, string(out))*/

		next.ServeHTTP(w, r)
	})
}
