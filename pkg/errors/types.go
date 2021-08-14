package errors

import (
    "errors"
    "github.com/go-chi/render"
    "net/http"
)

var (
    ErrUserNotFound            = errors.New("user not found")
    ErrIncorrectPassword       = errors.New("incorrect password")
    ErrUserAlreadyExist        = errors.New("user already exist")
    ErrCurrencyNotFound        = errors.New("currency not found")
    ErrBankAccountAlreadyExist = errors.New("bank account already exist")
    ErrBankAccountNotFound     = errors.New("bank account not found")
    ErrSomethingWrong          = errors.New("something went wrong")
    ErrCurrencyMismatch        = errors.New("currency mismatch")
    ErrWalletNotFound          = errors.New("wallet not found")
    ErrMissingAuthHeader       = errors.New("missing authorization header")
    ErrInvalidAuthHeaderFormat = errors.New("invalid auth header format")
    ErrUnsupportedAuth         = errors.New("auth type not supported")
    ErrUnauthorized            = errors.New("unauthorized user")
)

// Error renderer type for handling all sorts of errors.
type Error struct {
    Err            error  `json:"-"`                                                                            // low-level runtime error
    HTTPStatusCode int    `json:"-"`                                                                            // http response status code
    ErrorText      string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
}

// Render implements the github.com/go-chi/render.Renderer interface for ErrResponse
func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
    render.Status(r, e.HTTPStatusCode)
    return nil
}

func Status(err error) int {
    switch err {
    case ErrUserNotFound:
        return http.StatusNotFound
    case ErrIncorrectPassword:
        return http.StatusUnauthorized
    case ErrUserAlreadyExist:
        return http.StatusForbidden
    case ErrCurrencyNotFound:
        return http.StatusNotFound
    case ErrBankAccountAlreadyExist:
        return http.StatusForbidden
    case ErrBankAccountNotFound:
        return http.StatusNotFound
    case ErrCurrencyMismatch:
        return http.StatusConflict
    case ErrWalletNotFound:
        return http.StatusNotFound
    case ErrMissingAuthHeader:
        return http.StatusUnauthorized
    case ErrInvalidAuthHeaderFormat:
        return http.StatusUnauthorized
    case ErrUnsupportedAuth:
        return http.StatusUnauthorized
    case ErrUnauthorized:
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}

func ErrResponse(err error) render.Renderer {
    return &Error{
        Err:            err,
        HTTPStatusCode: Status(err),
        ErrorText:      err.Error(),
    }
}

func ErrBadRequest(err error) render.Renderer {
    return &Error{
        Err:            err,
        HTTPStatusCode: http.StatusBadRequest,
        ErrorText:      err.Error(),
    }
}
