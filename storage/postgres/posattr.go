package postgres

import (
	"context"
	"fmt"
	"position_service/pkg/helper"
	"position_service/storage"

	pb "position_service/genproto/position_service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type posattrRepo struct {
	db *pgxpool.Pool
}

func NewPosAttrRepo(db *pgxpool.Pool) storage.PosAttrI {
	return &posattrRepo{
		db: db,
	}
}

func (r *posattrRepo) Create(ctx context.Context, entity *pb.CreatePosAttrRequest) (id string, err error) {
	query := `
		INSERT INTO position_attributes (
			id,
			attribute_id,
			position_id,
			value
		) 
		 VALUES ($1, $2, $3,$4)
	`

	id = uuid.NewString()

	_, err = r.db.Exec(
		ctx,
		query,
		id,
		entity.AttributeId,
		entity.PositionId,
		entity.Value,
	)

	if err != nil {
		return "", fmt.Errorf("error while inserting posattr err: %w", err)
	}

	return id, nil
}

func (r *posattrRepo) GetAll(ctx context.Context, req *pb.GetAllPosAttrRequest) (*pb.GetAllPosAttrResponse, error) {
	var (
		resp   pb.GetAllPosAttrResponse
		err    error
		filter string
		params = make(map[string]interface{})
	)

	if req.Search != "" {
		filter += " AND value ILIKE '%' || :search || '%' "
		params["search"] = req.Search
	}

	countQuery := `SELECT count(1) FROM position_attributes WHERE true ` + filter

	q, arr := helper.ReplaceQueryParams(countQuery, params)
	err = r.db.QueryRow(ctx, q, arr...).Scan(
		&resp.Count,
	)

	if err != nil {
		return nil, fmt.Errorf("error while scanning count %w", err)
	}

	query := `SELECT
				id,
				attribute_id,
				position_id,
				value
			FROM position_attributes
			WHERE true` + filter

	query += " LIMIT :limit OFFSET :offset"
	params["limit"] = req.Limit
	params["offset"] = req.Offset

	q, arr = helper.ReplaceQueryParams(query, params)
	rows, err := r.db.Query(ctx, q, arr...)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows %w", err)
	}

	for rows.Next() {
		var posattr pb.PosAttr

		err = rows.Scan(
			&posattr.Id,
			&posattr.AttributeId,
			&posattr.PositionId,
			&posattr.Value,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning posattr err: %w", err)
		}

		resp.PosAttrs = append(resp.PosAttrs, &posattr)
	}

	return &resp, nil
}

func (r *posattrRepo) GetById(ctx context.Context, req *pb.GetByIdPosAttrRequest) (*pb.GetByIdPosAttrResponse, error) {
	var resp pb.GetByIdPosAttrResponse
	query := `
		SELECT * FROM position_attributes WHERE id=$1
	`
	rows, err := r.db.Query(ctx, query, req.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows %w", err)
	}

	for rows.Next() {
		var posattr pb.GetByIdPosAttrResponse

		err = rows.Scan(
			&posattr.Id,
			&posattr.AttributeId,
			&posattr.PositionId,
			&posattr.Value,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning posattr err: %w", err)
		}

		resp.Id = posattr.Id
		resp.AttributeId = posattr.AttributeId
		resp.PositionId = posattr.PositionId
		resp.Value = posattr.Value
	}

	return &resp, nil

}

func (r *posattrRepo) Update(ctx context.Context, entity *pb.UpdatePosAttrRequest) (*pb.UpdatePosAttrResponse, error) {
	var tempr pb.UpdatePosAttrResponse
	query := `
		UPDATE position_attributes SET attribute_id=$1 position_id=$2 value=$3 WHERE id=$4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.AttributeId,
		entity.PositionId,
		entity.Value,
		entity.Id,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while updating posattr err: %w", err)
	}

	return &tempr, nil
}

func (r *posattrRepo) Delete(ctx context.Context, entity *pb.DeletePosAttrRequest) (*pb.DeletePosAttrResponse, error) {
	var tempr pb.DeletePosAttrResponse
	query := `
		DELETE FROM position_attributes WHERE id=$1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.Id,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while deleting posattr err: %w", err)
	}

	return &tempr, nil
}
