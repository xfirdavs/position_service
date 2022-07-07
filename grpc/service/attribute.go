package service

import (
	"context"
	"position_service/config"
	pb "position_service/genproto/position_service"
	"position_service/pkg/logger"
	"position_service/storage"
)

type attributeService struct {
	cfg  config.Config
	log  logger.LoggerI
	strg storage.StorageI
	pb.UnimplementedAttributeServiceServer
}

func NewAttributeService(cfg config.Config, log logger.LoggerI, strg storage.StorageI) *attributeService {
	return &attributeService{
		cfg:  cfg,
		log:  log,
		strg: strg,
	}
}

func (s *attributeService) Create(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.Attribute, error) {
	id, err := s.strg.Attribute().Create(ctx, req)
	if err != nil {
		s.log.Error("CreateAttribute", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.Attribute{
		Id:   id,
		Name: req.Name,
		Type: req.Type,
	}, nil
}

func (s *attributeService) GetAll(ctx context.Context, req *pb.GetAllAttributeRequest) (*pb.GetAllAttributeResponse, error) {
	resp, err := s.strg.Attribute().GetAll(ctx, req)
	if err != nil {
		s.log.Error("GetAllAttribute", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *attributeService) GetById(ctx context.Context, req *pb.GetByIdAttributeRequest) (*pb.GetByIdAttributeResponse, error) {
	resp, err := s.strg.Attribute().GetById(ctx, req)
	if err != nil {
		s.log.Error("GetByIdAttribute", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (s *attributeService) Update(ctx context.Context, req *pb.UpdateAttributeRequest) (*pb.UpdateAttributeResponse, error) {
	_, err := s.strg.Attribute().Update(ctx, req)
	if err != nil {
		s.log.Error("UpdateAttribute", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.UpdateAttributeResponse{
		Id:   req.Id,
		Name: req.Name,
		Type: req.Type,
	}, nil
}

func (s *attributeService) Delete(ctx context.Context, req *pb.DeleteAttributeRequest) (*pb.DeleteAttributeResponse, error) {
	_, err := s.strg.Attribute().Delete(ctx, req)
	if err != nil {
		s.log.Error("DeleteAttribute", logger.Any("req", req), logger.Error(err))
		return nil, err
	}

	return &pb.DeleteAttributeResponse{
		Id: req.Id,
	}, nil
}
