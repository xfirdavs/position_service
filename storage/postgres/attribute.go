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

type attributeRepo struct {
	db *pgxpool.Pool
}

func NewAttributeRepo(db *pgxpool.Pool) storage.AttributeI {
	return &attributeRepo{
		db: db,
	}
}

func (r *attributeRepo) Create(ctx context.Context, entity *pb.CreateAttributeRequest) (id string, err error) {
	query := `
		INSERT INTO attribute (
			id,
			name,
			type
		) 
		 VALUES ($1, $2, $3)
	`

	id = uuid.NewString()

	_, err = r.db.Exec(
		ctx,
		query,
		id,
		entity.Name,
		entity.Type,
	)

	if err != nil {
		return "", fmt.Errorf("error while inserting attribute err: %w", err)
	}

	return id, nil
}

func (r *attributeRepo) GetAll(ctx context.Context, req *pb.GetAllAttributeRequest) (*pb.GetAllAttributeResponse, error) {
	var (
		resp   pb.GetAllAttributeResponse
		err    error
		filter string
		params = make(map[string]interface{})
	)

	if req.Search != "" {
		filter += " AND name ILIKE '%' || :search || '%' "
		params["search"] = req.Search
	}

	countQuery := `SELECT count(1) FROM attribute WHERE true ` + filter

	q, arr := helper.ReplaceQueryParams(countQuery, params)
	err = r.db.QueryRow(ctx, q, arr...).Scan(
		&resp.Count,
	)

	if err != nil {
		return nil, fmt.Errorf("error while scanning count %w", err)
	}

	query := `SELECT
				id,
				name,
				type
			FROM attribute
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
		var attribute pb.Attribute

		err = rows.Scan(
			&attribute.Id,
			&attribute.Name,
			&attribute.Type,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning attribute err: %w", err)
		}

		resp.Attributes = append(resp.Attributes, &attribute)
	}

	return &resp, nil
}

func (r *attributeRepo) GetById(ctx context.Context, req *pb.GetByIdAttributeRequest) (*pb.GetByIdAttributeResponse, error) {
	var resp pb.GetByIdAttributeResponse
	query := `
		SELECT * FROM attribute WHERE id=$1
	`
	rows, err := r.db.Query(ctx, query, req.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows %w", err)
	}

	for rows.Next() {
		var attribute pb.GetByIdAttributeResponse

		err = rows.Scan(
			&attribute.Id,
			&attribute.Name,
			&attribute.Type,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning attribute err: %w", err)
		}

		resp.Id = attribute.Id
		resp.Name = attribute.Name
		resp.Type = attribute.Type
	}

	return &resp, nil

}

func (r *attributeRepo) Update(ctx context.Context, entity *pb.UpdateAttributeRequest) (*pb.UpdateAttributeResponse, error) {
	var tempr pb.UpdateAttributeResponse
	query := `
		UPDATE attribute SET name=$1, type=$2 WHERE id=$3
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.Name,
		entity.Type,
		entity.Id,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while updating attribute err: %w", err)
	}

	return &tempr, nil
}

func (r *attributeRepo) Delete(ctx context.Context, entity *pb.DeleteAttributeRequest) (*pb.DeleteAttributeResponse, error) {
	var tempr pb.DeleteAttributeResponse
	query := `
		DELETE FROM attribute WHERE id=$1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.Id,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while deleting attribute err: %w", err)
	}

	return &tempr, nil
}
