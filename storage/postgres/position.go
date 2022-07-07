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

type positionRepo struct {
	db *pgxpool.Pool
}

func NewPositionRepo(db *pgxpool.Pool) storage.PositionI {
	return &positionRepo{
		db: db,
	}
}

func (r *positionRepo) Create(ctx context.Context, entity *pb.CreatePositionRequest) (id string, err error) {
	query := `
		INSERT INTO position (
			id,
			name,
			profession_id,
			company_id
		) 
		 VALUES ($1, $2, $3,$4)
	`

	id = uuid.NewString()

	_, err = r.db.Exec(
		ctx,
		query,
		id,
		entity.Name,
		entity.ProfessionId,
		entity.CompanyId,
	)

	if err != nil {
		return "", fmt.Errorf("error while inserting position err: %w", err)
	}

	return id, nil
}

func (r *positionRepo) GetAll(ctx context.Context, req *pb.GetAllPositionRequest) (*pb.GetAllPositionResponse, error) {
	var (
		resp   pb.GetAllPositionResponse
		err    error
		filter string
		params = make(map[string]interface{})
	)

	if req.Search != "" {
		filter += " AND name ILIKE '%' || :search || '%' "
		params["search"] = req.Search
	}

	countQuery := `SELECT count(1) FROM position WHERE true ` + filter

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
				profession_id,
				company_id
			FROM position
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
		var position pb.Position

		err = rows.Scan(
			&position.Id,
			&position.Name,
			&position.ProfessionId,
			&position.CompanyId,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning position err: %w", err)
		}

		resp.Positions = append(resp.Positions, &position)
	}

	return &resp, nil
}

func (r *positionRepo) GetById(ctx context.Context, req *pb.GetByIdPositionRequest) (*pb.GetByIdPositionResponse, error) {
	var resp pb.GetByIdPositionResponse
	query := `
		SELECT * FROM position WHERE id=$1
	`
	rows, err := r.db.Query(ctx, query, req.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows %w", err)
	}

	for rows.Next() {
		var position pb.GetByIdPositionResponse

		err = rows.Scan(
			&position.Id,
			&position.Name,
			&position.ProfessionId,
			&position.CompanyId,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning position err: %w", err)
		}

		resp.Id = position.Id
		resp.Name = position.Name
		resp.ProfessionId = position.ProfessionId
		resp.CompanyId = position.CompanyId
	}

	return &resp, nil

}

func (r *positionRepo) Update(ctx context.Context, entity *pb.UpdatePositionRequest) (*pb.UpdatePositionResponse, error) {
	var tempr pb.UpdatePositionResponse
	query := `
		UPDATE position SET name=$1 type=$2 WHERE id=$3
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.Name,
		entity.Id,
		entity.ProfessionId,
		entity.CompanyId,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while updating position err: %w", err)
	}

	return &tempr, nil
}

func (r *positionRepo) Delete(ctx context.Context, entity *pb.DeletePositionRequest) (*pb.DeletePositionResponse, error) {
	var tempr pb.DeletePositionResponse
	query := `
		DELETE FROM position WHERE id=$1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		entity.Id,
	)

	if err != nil {
		return &tempr, fmt.Errorf("error while deleting position err: %w", err)
	}

	return &tempr, nil
}
