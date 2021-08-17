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
    "strconv"
)

type WalletResource interface {
    Pay(w http.ResponseWriter, r *http.Request)
    Get(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router)
}

type walletResource struct {
    walletSvc service.WalletSvc
}

func NewWalletResource(walletSvc service.WalletSvc) WalletResource {
    return &walletResource{
        walletSvc: walletSvc,
    }
}

func (wr *walletResource) RegisterRoutes(r chi.Router) {
    r.Get("/wallets/{walletID}", wr.Get)
    r.Post("/wallets/pay", wr.Pay)
}

func (wr *walletResource) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    walletID := chi.URLParam(r, "walletID")

    id, err := strconv.Atoi(walletID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    res, err := wr.walletSvc.GetWalletById(ctx, int64(id))
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (wr *walletResource) Pay(w http.ResponseWriter, r *http.Request) {
    var req dto.TransferMoneyDto
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

    res, err := wr.walletSvc.Pay(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}
