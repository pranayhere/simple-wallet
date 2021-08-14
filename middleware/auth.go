package middleware

import (
    "context"
    "github.com/go-chi/render"
    "github.com/pranayhere/simple-wallet/pkg/constant"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/token"
    "net/http"
    "strings"
)

func Auth(tokenMaker token.Maker) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authorizationHeader := r.Header.Get(constant.AuthorizationHeaderKey)
            if len(authorizationHeader) == 0 {
                _ = render.Render(w, r, errors.ErrResponse(errors.ErrMissingAuthHeader))
                return
            }

            fields := strings.Fields(authorizationHeader)
            if len(fields) < 2 {
                _ = render.Render(w, r, errors.ErrResponse(errors.ErrInvalidAuthHeaderFormat))
                return
            }

            authorizationType := strings.ToLower(fields[0])
            if authorizationType != constant.AuthorizationTypeBearer {
                _ = render.Render(w, r, errors.ErrResponse(errors.ErrUnsupportedAuth))
                return
            }

            accessToken := fields[1]
            payload, err := tokenMaker.VerifyToken(accessToken)
            if err != nil {
                _ = render.Render(w, r, errors.ErrResponse(errors.ErrUnauthorized))
                return
            }

            ctx := context.WithValue(r.Context(), constant.AuthorizationPayloadKey, payload)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
