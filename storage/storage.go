package storage

import (
	"context"
	"errors"
	pb "position_service/genproto/position_service"
)

var ErrorTheSameId = errors.New("cannot use the same uuid for 'id' and 'parent_id' fields")
var ErrorProjectId = errors.New("not valid 'project_id'")

type StorageI interface {
	Profession() ProfessionI
	Attribute() AttributeI
	Position() PositionI
	PosAttr() PosAttrI
}

type ProfessionI interface {
	Create(ctx context.Context, entity *pb.CreateProfessionRequest) (id string, err error)
	GetAll(ctx context.Context, req *pb.GetAllProfessionRequest) (*pb.GetAllProfessionResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdProfessionRequest) (*pb.GetByIdProfessionResponse, error)
	Update(ctx context.Context, entity *pb.UpdateProfessionRequest) (*pb.UpdateProfessionResponse, error)
	Delete(ctx context.Context, entity *pb.DeleteProfessionRequest) (*pb.DeleteProfessionResponse, error)
}

type AttributeI interface {
	Create(ctx context.Context, entity *pb.CreateAttributeRequest) (id string, err error)
	GetAll(ctx context.Context, req *pb.GetAllAttributeRequest) (*pb.GetAllAttributeResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdAttributeRequest) (*pb.GetByIdAttributeResponse, error)
	Update(ctx context.Context, entity *pb.UpdateAttributeRequest) (*pb.UpdateAttributeResponse, error)
	Delete(ctx context.Context, entity *pb.DeleteAttributeRequest) (*pb.DeleteAttributeResponse, error)
}

type PositionI interface {
	Create(ctx context.Context, entity *pb.CreatePositionRequest) (id string, err error)
	GetAll(ctx context.Context, req *pb.GetAllPositionRequest) (*pb.GetAllPositionResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdPositionRequest) (*pb.GetByIdPositionResponse, error)
	Update(ctx context.Context, entity *pb.UpdatePositionRequest) (*pb.UpdatePositionResponse, error)
	Delete(ctx context.Context, entity *pb.DeletePositionRequest) (*pb.DeletePositionResponse, error)
}

type PosAttrI interface {
	Create(ctx context.Context, entity *pb.CreatePosAttrRequest) (id string, err error)
	GetAll(ctx context.Context, req *pb.GetAllPosAttrRequest) (*pb.GetAllPosAttrResponse, error)
	GetById(ctx context.Context, req *pb.GetByIdPosAttrRequest) (*pb.GetByIdPosAttrResponse, error)
	Update(ctx context.Context, entity *pb.UpdatePosAttrRequest) (*pb.UpdatePosAttrResponse, error)
	Delete(ctx context.Context, entity *pb.DeletePosAttrRequest) (*pb.DeletePosAttrResponse, error)
}
