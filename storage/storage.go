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
}

type ProfessionI interface {
	Create(ctx context.Context, entity *pb.CreateProfessionRequest) (id string, err error)
	GetAll(ctx context.Context, req *pb.GetAllProfessionRequest) (*pb.GetAllProfessionResponse, error)
}
