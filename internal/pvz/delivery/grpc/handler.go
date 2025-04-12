package grpc

import (
	"context"

	pb "github.com/0x0FACED/pvz-avito/internal/pvz/delivery/grpc/v1"
	pvz_domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZService interface {
	ListAllPVZs(ctx context.Context) ([]*pvz_domain.PVZ, error)
}

type GRPCHandler struct {
	pb.UnimplementedPVZServiceServer
	svc PVZService
}

func NewGRPCHandler(svc PVZService) *GRPCHandler {
	return &GRPCHandler{
		svc: svc,
	}
}

func (h *GRPCHandler) GetPVZList(ctx context.Context, req *pb.GetPVZListRequest) (*pb.GetPVZListResponse, error) {
	pvzs, err := h.svc.ListAllPVZs(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetPVZListResponse{}
	for _, p := range pvzs {
		resp.Pvzs = append(resp.Pvzs, &pb.PVZ{
			Id:               *p.ID,
			RegistrationDate: timestamppb.New(*p.RegistrationDate),
			City:             p.City.String(),
		})
	}

	return resp, nil
}
