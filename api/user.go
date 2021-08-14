package api

import (
    "encoding/json"
    "github.com/go-chi/chi"
    "github.com/go-chi/render"
    "github.com/go-playground/validator/v10"
    "github.com/pranayhere/simple-wallet/dto"
    types "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/sirupsen/logrus"
    "net/http"
)

type UserResource interface {
    Create(w http.ResponseWriter, r *http.Request)
    Login(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router)
}

type userResource struct {
    userSvc service.UserSvc
}

func NewUserResource(userSvc service.UserSvc) UserResource {
    return &userResource{
        userSvc: userSvc,
    }
}

func (u *userResource) RegisterRoutes(r chi.Router) {
    r.Post("/users", u.Create)
    r.Post("/users/login", u.Login)
}

func (u *userResource) Create(w http.ResponseWriter, r *http.Request) {
    logrus.Println("log create user")
    var req dto.CreateUserDto
    ctx := r.Context()

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }
    defer r.Body.Close()

    if err := validator.New().Struct(req); err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    user, err := u.userSvc.CreateUser(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, user)
}

func (u *userResource) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginCredentialsDto
    ctx := r.Context()

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }
    defer r.Body.Close()

    if err := validator.New().Struct(req); err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    loggedInUser, err := u.userSvc.LoginUser(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, loggedInUser)
}
