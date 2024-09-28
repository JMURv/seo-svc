package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/par-pro-seo/api/pb"
	ctrl "github.com/JMURv/par-pro-seo/internal/controller"
	hdl "github.com/JMURv/par-pro-seo/internal/handler"
	metrics "github.com/JMURv/par-pro-seo/internal/metrics/prometheus"
	"github.com/JMURv/par-pro-seo/internal/validation"
	utils "github.com/JMURv/par-pro-seo/pkg/utils/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) GetSEO(ctx context.Context, req *pb.GetReq) (*pb.SEOMsg, error) {
	const op = "seo.GetSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	res, err := h.ctrl.GetSEO(ctx, req.Name, req.Pk)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return utils.ModelToProto(res), nil
}

func (h *Handler) CreateSEO(ctx context.Context, req *pb.CreateAndUpdateReq) (*pb.Empty, error) {
	const op = "seo.CreateSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	obj := utils.ProtoToModel(req.Seo)
	if err := validation.ValidateSEO(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	err := h.ctrl.CreateSEO(ctx, req.Name, req.Pk, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.Empty{}, nil
}

func (h *Handler) UpdateSEO(ctx context.Context, req *pb.CreateAndUpdateReq) (*pb.Empty, error) {
	const op = "seo.UpdateSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	obj := utils.ProtoToModel(req.Seo)
	if err := validation.ValidateSEO(obj); err != nil {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, err.Error())
	}

	err := h.ctrl.UpdateSEO(ctx, req.Name, req.Pk, obj)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.Empty{}, nil
}

func (h *Handler) DeleteSEO(ctx context.Context, req *pb.GetReq) (*pb.Empty, error) {
	const op = "seo.DeleteSEO.hdl"
	s, c := time.Now(), codes.OK
	span := opentracing.GlobalTracer().StartSpan(op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), int(c), op)
	}()

	if req == nil || req.Name == "" || req.Pk == "" {
		c = codes.InvalidArgument
		return nil, status.Errorf(c, hdl.ErrDecodeRequest.Error())
	}

	err := h.ctrl.DeleteSEO(ctx, req.Name, req.Pk)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = codes.NotFound
		return nil, status.Errorf(c, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		c = codes.Internal
		return nil, status.Errorf(c, hdl.ErrInternal.Error())
	}
	return &pb.Empty{}, nil
}
