package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/GirishBhutiya/gOpenfiatServer/app/handler"
	"github.com/GirishBhutiya/gOpenfiatServer/app/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

var maker token.Maker

func InitAuthTokenMaker(mk *token.Maker) {
	maker = *mk
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			handler.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			handler.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			handler.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]

		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			handler.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}
		//log.Println("middleware", payload.UserId)
		ctx := context.WithValue(r.Context(), handler.UserIDKey, payload.UserId)
		r = r.WithContext(ctx)
		/*out, err := json.Marshal(payload)
		if err != nil {
			app.errorJSON(w, err, http.StatusUnauthorized)
			return
		}
		r.Header.Set(authorizationPayloadKey, string(out))*/

		next.ServeHTTP(w, r)
	})
}
