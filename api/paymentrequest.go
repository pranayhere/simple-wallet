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

type PaymentRequestResource interface {
    Charge(w http.ResponseWriter, r *http.Request)
    Approve(w http.ResponseWriter, r *http.Request)
    RegisterRoutes(r chi.Router)
}

type paymentRequestResource struct {
    payReqSvc service.PaymentRequestSvc
}

func NewPaymentRequestResource(payReqSvc service.PaymentRequestSvc) PaymentRequestResource {
    return &paymentRequestResource{
        payReqSvc: payReqSvc,
    }
}

func (p *paymentRequestResource) RegisterRoutes(r chi.Router) {
    r.Post("/payment-req", p.Charge)
    r.Get("/payment-req", p.List)
    r.Patch("/payment-req/{payReqID}/approve", p.Approve)
    r.Patch("/payment-req/{payReqID}/refuse", p.Refuse)
}

func (p *paymentRequestResource) Charge(w http.ResponseWriter, r *http.Request) {
    var req dto.PaymentRequestDto
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

    res, err := p.payReqSvc.Create(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (p *paymentRequestResource) Approve(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    payReqID := chi.URLParam(r, "payReqID")

    id, err := strconv.Atoi(payReqID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    res, err := p.payReqSvc.Approve(ctx, int64(id))
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (p *paymentRequestResource) Refuse(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    payReqID := chi.URLParam(r, "payReqID")

    id, err := strconv.Atoi(payReqID)
    if err != nil {
        _ = render.Render(w, r, types.ErrBadRequest(err))
        return
    }

    res, err := p.payReqSvc.Refuse(ctx, int64(id))
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}

func (p *paymentRequestResource) List(w http.ResponseWriter, r *http.Request) {
    var req dto.ListPaymentRequestsDto
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

    res, err := p.payReqSvc.List(ctx, req)
    if err != nil {
        _ = render.Render(w, r, types.ErrResponse(err))
        return
    }

    render.JSON(w, r, res)
}
