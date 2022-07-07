package service

import (
	"context"
	"position_service/config"
	pb "position_service/genproto/position_service"
	"position_service/pkg/logger"
	"position_service/storage"
)

type positionService struct {
	cfg  config.Config
	log  logger.LoggerI
	strg storage.StorageI
	pb.UnimplementedPositionServiceServer
}

func NewPositionService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *positionService {
	return &positionService{
		cfg:  cfg,
		log:  log,
		strg: strg,
	}
}

func (s *positionService) Create(ctx context.Context, req *pb.CreatePositionRequest) (*pb.Position, error) {
	id, err := s.strg.Position().Create(ctx, req)
	if err != nil {
		s.log.Error("CreatePosition", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.Position{
		Id:           id,
		Name:         req.Name,
		ProfessionId: req.ProfessionId,
		CompanyId:    req.CompanyId,
	}, nil
}

func (s *positionService) GetAll(ctx context.Context, req *pb.GetAllPositionRequest) (*pb.GetAllPositionResponse, error) {
	resp, err := s.strg.Position().GetAll(ctx, req)
	if err != nil {
		s.log.Error("GetAllPosition", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *positionService) GetById(ctx context.Context, req *pb.GetByIdPositionRequest) (*pb.GetByIdPositionResponse, error) {
	resp, err := s.strg.Position().GetById(ctx, req)
	if err != nil {
		s.log.Error("GetByIdPosition", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *positionService) Update(ctx context.Context, req *pb.UpdatePositionRequest) (*pb.UpdatePositionResponse, error) {
	_, err := s.strg.Position().Update(ctx, req)
	if err != nil {
		s.log.Error("UpdatePosition", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.UpdatePositionResponse{
		Id:           req.Id,
		Name:         req.Name,
		ProfessionId: req.ProfessionId,
		CompanyId:    req.CompanyId,
	}, nil
}

func (s *positionService) Delete(ctx context.Context, req *pb.DeletePositionRequest) (*pb.DeletePositionResponse, error) {
	_, err := s.strg.Position().Delete(ctx, req)
	if err != nil {
		s.log.Error("DeletePosition", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.DeletePositionResponse{
		Id: req.Id,
	}, nil
}
