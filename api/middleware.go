package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mbeka02/bank/internal/auth"
)

const (
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func denyAccess(w http.ResponseWriter, msg string) {
	JSONResponse(w, http.StatusUnauthorized, msg)
}
func getAuthPayload(ctx context.Context) (*auth.Payload, error) {
	payload, ok := ctx.Value(authorizationPayloadKey).(*auth.Payload)

	if !ok {
		return nil, APIError{
			message:    "internal server error",
			statusCode: http.StatusInternalServerError,
		}
	}

	return payload, nil
}

func authMiddleware(next http.Handler, tokenMaker auth.Maker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		next.ServeHTTP(w, r)
	})
}
