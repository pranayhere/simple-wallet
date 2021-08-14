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

type BankAccountResource interface {
    Create(w http.ResponseWriter, r *http.Request)
    VerificationSuccess(w http.ResponseWriter, r *http.Request)
    VerificationFailed(w http.ResponseWriter, r *http.Request)
    Get(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router)
}

type bankAccountResource struct {
    bankAcctSvc service.BankAccountSvc
}

func NewBankAccountResource(bankAcctSvc service.BankAccountSvc) BankAccountResource {
    return &bankAccountResource{
        bankAcctSvc: bankAcctSvc,
    }
}

func (b *bankAccountResource) RegisterRoutes(r chi.Router) {
    r.Get("/bank-accounts/{bankAcctID}", b.Get)
    r.Post("/bank-accounts", b.Create)
    r.Patch("/bank-accounts/{bankAcctID}/verification-success", b.VerificationSuccess)
    r.Patch("/bank-accounts/{bankAcctID}/verification-failed", b.VerificationFailed)
}

func (b *bankAccountResource) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateBankAccountDto
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

    res, err := b.bankAcctSvc.CreateBankAccount(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (b *bankAccountResource) VerificationSuccess(w http.ResponseWriter, r *http.Request) {
    var req dto.BankAccountVerificationDto
    ctx := r.Context()

    bankAcctID := chi.URLParam(r, "bankAcctID")
    id, err := strconv.Atoi(bankAcctID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    req = dto.BankAccountVerificationDto{
        BankAccountID: int64(id),
    }

    res, err := b.bankAcctSvc.VerificationSuccess(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (b *bankAccountResource) VerificationFailed(w http.ResponseWriter, r *http.Request) {
    var req dto.BankAccountVerificationDto
    ctx := r.Context()

    bankAcctID := chi.URLParam(r, "bankAcctID")
    id, err := strconv.Atoi(bankAcctID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    req = dto.BankAccountVerificationDto{
        BankAccountID: int64(id),
    }

    res, err := b.bankAcctSvc.VerificationFailed(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (b *bankAccountResource) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    bankAcctID := chi.URLParam(r, "bankAcctID")

    id, err := strconv.Atoi(bankAcctID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    res, err := b.bankAcctSvc.GetBankAccount(ctx, int64(id))
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}
