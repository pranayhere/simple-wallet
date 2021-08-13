package api

import (
    "encoding/json"
    "github.com/go-chi/chi"
    "github.com/go-chi/render"
    "github.com/go-playground/validator/v10"
    "github.com/pranayhere/simple-wallet/dto"
    types "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/service"
    "net/http"
)

type CurrencyResource interface {
    Get(w http.ResponseWriter, r *http.Request)
    Create(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router) http.Handler
}

type currencyResource struct {
    currencySvc service.CurrencySvc
}

func (s *currencyResource) RegisterRoutes(r chi.Router) http.Handler {
    r.Get("/{currencyCode}", s.Get)
    r.Post("/", s.Create)

    return r
}

func NewCurrencyResource(currencySvc service.CurrencySvc) CurrencyResource {
    return &currencyResource{
        currencySvc: currencySvc,
    }
}

func (s *currencyResource) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    code := chi.URLParam(r, "currencyCode")

    res, err := s.currencySvc.GetCurrency(ctx, code)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (s *currencyResource) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CurrencyDto
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

    res, err := s.currencySvc.CreateCurrency(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}