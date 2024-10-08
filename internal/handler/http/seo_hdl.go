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

func RegisterSEORoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/seo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateSEO(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/seo/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetSEO(w, r)
		case http.MethodPut:
			h.UpdateSEO(w, r)
		case http.MethodDelete:
			h.DeleteSEO(w, r)
		default:
			utils.ErrResponse(w, http.StatusMethodNotAllowed, hdl.ErrMethodNotAllowed)
		}
	})
}

func (h *Handler) GetSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "seo.GetItemSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	res, err := h.ctrl.GetSEO(r.Context(), name, pk)
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

func (h *Handler) CreateSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusCreated
	const op = "seo.CreateSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &model.SEO{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSEO(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.CreateSEO(r.Context(), req)
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

func (h *Handler) UpdateSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "seo.UpdateSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &model.SEO{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSEO(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdateSEO(r.Context(), req)
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

func (h *Handler) DeleteSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusNoContent
	const op = "seo.DeleteSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request, missing name or pk",
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	err := h.ctrl.DeleteSEO(r.Context(), name, pk)
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
