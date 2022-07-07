package service

import (
	"context"
	"position_service/config"
	pb "position_service/genproto/position_service"
	"position_service/pkg/logger"
	"position_service/storage"
)

type posattrService struct {
	cfg  config.Config
	log  logger.LoggerI
	strg storage.StorageI
	pb.UnimplementedPosAttrServiceServer
}

func NewPosAttrService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *posattrService {
	return &posattrService{
		cfg:  cfg,
		log:  log,
		strg: strg,
	}
}

func (s *posattrService) Create(ctx context.Context, req *pb.CreatePosAttrRequest) (*pb.PosAttr, error) {
	id, err := s.strg.PosAttr().Create(ctx, req)
	if err != nil {
		s.log.Error("CreatePosAttr", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.PosAttr{
		Id:          id,
		AttributeId: req.AttributeId,
		PositionId:  req.PositionId,
		Value:       req.Value,
	}, nil
}

func (s *posattrService) GetAll(ctx context.Context, req *pb.GetAllPosAttrRequest) (*pb.GetAllPosAttrResponse, error) {
	resp, err := s.strg.PosAttr().GetAll(ctx, req)
	if err != nil {
		s.log.Error("GetAllPosAttr", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *posattrService) GetById(ctx context.Context, req *pb.GetByIdPosAttrRequest) (*pb.GetByIdPosAttrResponse, error) {
	resp, err := s.strg.PosAttr().GetById(ctx, req)
	if err != nil {
		s.log.Error("GetByIdPosAttr", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *posattrService) Update(ctx context.Context, req *pb.UpdatePosAttrRequest) (*pb.UpdatePosAttrResponse, error) {
	_, err := s.strg.PosAttr().Update(ctx, req)
	if err != nil {
		s.log.Error("UpdatePosAttr", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.UpdatePosAttrResponse{
		Id:          req.Id,
		AttributeId: req.AttributeId,
		PositionId:  req.PositionId,
		Value:       req.Value,
	}, nil
}

func (s *posattrService) Delete(ctx context.Context, req *pb.DeletePosAttrRequest) (*pb.DeletePosAttrResponse, error) {
	_, err := s.strg.PosAttr().Delete(ctx, req)
	if err != nil {
		s.log.Error("DeletePosAttr", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.DeletePosAttrResponse{
		Id: req.Id,
	}, nil
}
