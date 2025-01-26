package postgres

import (
	"context"
	"main/internal/domain/models"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	db *pgxpool.Pool
}

func NewRequestRepository(db *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Create(ctx context.Context, request models.Request) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := squirrel.Insert("requests").
		Columns("id", "name", "phone", "email", "car_info", "date").
		Values(request.ID, request.Name, request.Phone, request.Email, request.CarInfo, request.Date).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Request, error) {
	query := squirrel.Select("id", "name", "phone", "email", "car_info", "date").
		From("requests").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return models.Request{}, err
	}

	var request models.Request
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&request.ID,
		&request.Name,
		&request.Phone,
		&request.Email,
		&request.CarInfo,
		&request.Date,
	)
	if err != nil {
		return models.Request{}, err
	}

	return request, nil
}

func (r *RequestRepository) List(ctx context.Context, page, pageSize int) ([]models.Request, int, error) {
	query := squirrel.Select(
		"id", "name", "phone", "email", "car_info", "COUNT(*) OVER() AS total_records").
		From("requests").
		PlaceholderFormat(squirrel.Dollar).
		Limit(uint64(pageSize)).
		Offset(uint64((page - 1) * pageSize))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var requests []models.Request
	var totalRecords int
	for rows.Next() {
		var request models.Request

		err := rows.Scan(
			&request.ID, &request.Name, &request.Phone, &request.Email, &request.CarInfo, &totalRecords,
		)
		if err != nil {
			return nil, 0, err
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	totalPages := (totalRecords + pageSize - 1) / pageSize

	return requests, totalPages, nil
}
