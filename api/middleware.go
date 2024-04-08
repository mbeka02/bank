package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mbeka02/bank/internal/auth"
)

func denyAccess(w http.ResponseWriter, msg string) {
	JSONResponse(w, http.StatusUnauthorized, msg)
}

const (
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(handlerFunc http.HandlerFunc, tokenMaker auth.Maker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			denyAccess(w, "authorization header is not set")
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			denyAccess(w, "malformed key")
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			denyAccess(w, "this authorization type is not supported")
			return
		}
		tokenString := fields[1]
		payload, err := tokenMaker.ValidateToken(tokenString)
		if err != nil {
			fmt.Println(err)
			denyAccess(w, "authentication failed")
			return
		}
		//refactor this
		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, authorizationPayloadKey, payload))
		*r = *req
		handlerFunc(w, r)
	}
}
