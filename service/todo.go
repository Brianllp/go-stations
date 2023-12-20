package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// ref: https://qiita.com/__m/items/b0cc9a4818297f5becd6
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	// ref: https://zenn.dev/chillout2san/articles/34aa952747880c#crud%E5%87%A6%E7%90%86%E3%81%AEread%E5%87%A6%E7%90%86%E3%82%92%E3%81%97%E3%81%9F%E3%81%84%E5%A0%B4%E5%90%88(%E4%B8%80%E3%81%A4%E3%81%AE%E3%83%AC%E3%82%B3%E3%83%BC%E3%83%89%E3%82%92%E8%BF%94%E3%81%99)
	row := s.db.QueryRowContext(ctx, confirm, id)
	if err := row.Scan(&subject, &description, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	} else if count == 0 {
		return nil, &model.ErrNotFound{}
	}

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	row := s.db.QueryRowContext(ctx, confirm, id)
	if err := row.Scan(&subject, &description, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
