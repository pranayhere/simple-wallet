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
    SendMoney(w http.ResponseWriter, r *http.Request)
    Deposit(w http.ResponseWriter, r *http.Request)
    Withdraw(w http.ResponseWriter, r *http.Request)
    Get(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router) http.Handler
}

type walletResource struct {
    walletSvc service.WalletSvc
}

func NewWalletResource(walletSvc service.WalletSvc) WalletResource {
    return &walletResource{
        walletSvc: walletSvc,
    }
}

func (wr *walletResource) RegisterRoutes(r chi.Router) http.Handler {
    r.Get("/{walletID}", wr.Get)
    r.Post("/send", wr.SendMoney)
    r.Post("/deposit", wr.Deposit)
    r.Post("/withdraw", wr.Withdraw)

    return r
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

func (wr *walletResource) SendMoney(w http.ResponseWriter, r *http.Request) {
    var req dto.SendMoneyDto
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

    res, err := wr.walletSvc.SendMoney(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (wr *walletResource) Deposit(w http.ResponseWriter, r *http.Request) {
    var req dto.DepositDto
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

    res, err := wr.walletSvc.Deposit(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (wr *walletResource) Withdraw(w http.ResponseWriter, r *http.Request) {
    var req dto.WithdrawDto
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

    res, err := wr.walletSvc.Withdraw(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}
