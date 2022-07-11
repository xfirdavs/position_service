package postgres

import (
	"context"
	"fmt"
	"position_service/pkg/helper"
	"position_service/storage"

	pb "position_service/genproto/position_service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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

func (r *positionRepo) Create(ctx context.Context, entity *pb.CreatePositionRequest) (id *pb.PositionId, err error) {
	positionQuery := `
		INSERT INTO position(
			id,
			name,
			profession_id,
			company_id
		) VALUES ($1, $2, $3, $4);

	`
	positionAttributesQuery := `
			INSERT INTO position_attributes(
				id,
				attribute_id,
				position_id,
				value
			)VALUES ($1,$2,$3,$4)
		`
	positionId := uuid.New()

	client, err := r.db.Acquire(context.TODO())
	transact, err := client.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer transact.Rollback(context.TODO())
	_, err = transact.Exec(
		ctx,
		positionQuery,
		positionId,
		entity.Name,
		entity.ProfessionId,
		entity.CompanyId,
	)
	if err != nil {
		return nil, err
	}
	positionAttributes := entity.PositionAttributes

	for k := range positionAttributes {
		attribute_id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		_, err = transact.Exec(
			ctx,
			positionAttributesQuery,
			attribute_id,
			positionAttributes[k].AttributeId,
			positionId,
			positionAttributes[k].Value,
		)
		if err != nil {
			return nil, err
		}
	}

	err = transact.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.PositionId{
		Id: positionId.String(),
	}, nil
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
		xpositionAttributes, err := r.GetPositionAttributesExample(ctx, position.Id)
		if err != nil {
			return nil, err
		}

		position.PositionAttributes = xpositionAttributes

		resp.Positions = append(resp.Positions, &position)
	}

	return &resp, nil
}

func (r *positionRepo) GetPositionAttributesExample(ctx context.Context, id string) ([]*pb.GetPositionAttributes, error) {
	var resp []pb.GetPositionAttributes
	query := `SELECT 
	pa.id, 
	pa.attribute_id,
	pa.position_id,
	pa.value,
	ap.name,
	ap.type 
	FROM position_attributes AS pa JOIN attribute ap ON ap.id=pa.attribute_id WHERE pa.position_id=$1
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting rows %w", err)
	}

	for rows.Next() {
		var position_attibutes pb.GetPositionAttributes

		err = rows.Scan(
			&position_attibutes.Id,
			&position_attibutes.AttributeId,
			&position_attibutes.PositionId,
			&position_attibutes.Value,
			&position_attibutes.AttributeName,
			&position_attibutes.AttributeType,
		)

		if err != nil {
			return nil, fmt.Errorf("error while scanning position err: %w", err)
		}

		resp = append(resp, position_attibutes)

	}
	return &resp, nil

}

func (r *positionRepo) GetById(ctx context.Context, req *pb.PositionId) (*pb.Position, error) {

	var resp pb.Position
	query := `
		SELECT id,name,profession_id,company_id FROM position WHERE id=$1
	`
	rows, err := r.db.Query(ctx, query, req.Id)
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

		resp.Id = position.Id
		resp.Name = position.Name
		resp.ProfessionId = position.ProfessionId
		resp.CompanyId = position.CompanyId
	}
	xpositionAttributes, err := r.GetPositionAttributesExample(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp.PositionAttributes = xpositionAttributes

	return &resp, nil

}

func (r *positionRepo) Update(ctx context.Context, entity *pb.UpdatePositionRequest) (*pb.PositionId, error) {
	var id pb.PositionId
	id.Id = entity.Id
	queryPosition := `
		UPDATE position SET name=$1, profession_id=$2, company_id=$3 WHERE id=$4
	`
	queryDeletePosAttr := `DELETE FROM position_attributes WHERE position_id=$1`
	queryPosAttr := `
		INSERT INTO position_attributes(
			id,
			attribute_id,
			position_id,
			value
		)VALUES ($1,$2,$3,$4);
	`

	client, err := r.db.Acquire(context.TODO())
	transact, err := client.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer transact.Rollback(context.TODO())
	_, err = transact.Exec(
		ctx,
		queryPosition,
		entity.Name,
		entity.ProfessionId,
		entity.CompanyId,
		id.Id,
	)
	if err != nil {
		return nil, err
	}

	_, err = transact.Exec(
		ctx,
		queryDeletePosAttr,
		entity.Id,
	)
	positionAttributes := entity.PositionAttributes
	for k := range positionAttributes {
		attribute_id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		_, err = transact.Exec(
			ctx,
			queryPosAttr,
			attribute_id,
			positionAttributes[k].AttributeId,
			entity.Id,
			positionAttributes[k].Value,
		)
		if err != nil {
			return nil, err
		}
	}

	err = transact.Commit(ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return &id, fmt.Errorf("error while updating position err: %w", err)
	}

	return &id, nil
}

func (r *positionRepo) Delete(ctx context.Context, entity *pb.PositionId) (*pb.PositionId, error) {
	var id pb.PositionId
	id.Id = entity.Id
	queryPositionAttributes := `
		DELETE FROM position_attributes WHERE position_id=$1
	`
	queryPosition := `
		DELETE FROM position WHERE id=$1
	`

	client, err := r.db.Acquire(ctx)
	transact, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer transact.Rollback(ctx)
	_, err = transact.Exec(
		ctx,
		queryPositionAttributes,
		id.Id,
	)
	if err != nil {
		return nil, err
	}

	_, err = transact.Exec(
		ctx,
		queryPosition,
		entity.Id,
	)

	err = transact.Commit(ctx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("error while deleting position err: %w", err)
	}

	return &id, nil
}
