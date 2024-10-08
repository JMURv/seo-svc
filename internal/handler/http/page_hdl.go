package http

import (
	"errors"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	hdl "github.com/JMURv/seo-svc/internal/handler"
	metrics "github.com/JMURv/seo-svc/internal/metrics/prometheus"
	"github.com/JMURv/seo-svc/internal/validation"
	"github.com/JMURv/seo-svc/pkg/model"
	utils "github.com/JMURv/seo-svc/pkg/utils/http"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegisterPageRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/page", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.ListPages(w, r)
		case http.MethodPost:
			h.CreatePage(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/page/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetPage(w, r)
		case http.MethodPut:
			h.UpdatePage(w, r)
		case http.MethodDelete:
			h.DeletePage(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})
}

func (h *Handler) ListPages(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "pages.ListPages.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	res, err := h.ctrl.ListPages(r.Context())
	if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "pages.GetPage.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op), zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	res, err := h.ctrl.GetPage(r.Context(), slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusCreated
	const op = "pages.CreatePage.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &model.Page{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request",
			zap.String("op", op), zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ValidatePage(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to validate obj",
			zap.String("op", op), zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.CreatePage(r.Context(), req)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		c = http.StatusConflict
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "pages.UpdatePage.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op), zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &model.Page{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	if err := validation.ValidatePage(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdatePage(r.Context(), slug, req)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) DeletePage(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "pages.DeletePage.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op),
			zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	err := h.ctrl.DeletePage(r.Context(), slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}
