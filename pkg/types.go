package types

import (
    "github.com/go-chi/render"
    "github.com/pranayhere/simple-wallet/common"
    "net/http"
)

// Error renderer type for handling all sorts of errors.
type Error struct {
    Err            error `json:"-"` // low-level runtime error
    HTTPStatusCode int   `json:"-"` // http response status code
    ErrorText  string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
}

// Render implements the github.com/go-chi/render.Renderer interface for ErrResponse
func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
    render.Status(r, e.HTTPStatusCode)
    return nil
}

func Status(err error) int {
    switch err {
    case common.ErrUserNotFound:
        return http.StatusNotFound
    case common.ErrIncorrectPassword:
        return http.StatusUnauthorized
    case common.ErrUserAlreadyExist:
        return http.StatusConflict
    case common.ErrCurrencyNotFound:
        return http.StatusNotFound
    case common.ErrBankAccountAlreadyExist:
        return http.StatusConflict
    case common.ErrBankAccountNotFound:
        return http.StatusNotFound
    case common.ErrCurrencyMismatch:
        return http.StatusConflict
    case common.ErrWalletNotFound:
        return http.StatusNotFound
    default:
        return http.StatusInternalServerError
    }
}

func ErrResponse(err error) render.Renderer {
    return &Error{
        Err: err,
        HTTPStatusCode: Status(err),
        ErrorText: err.Error(),
    }
}

func ErrValidation(err error) render.Renderer {
    return &Error{
        Err: err,
        HTTPStatusCode: http.StatusBadRequest,
        ErrorText: err.Error(),
    }
}