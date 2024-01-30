package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"note/internal/models"
	"time"
)

type noteStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *noteStorage {
	return &noteStorage{db: db}
}

func (ns *noteStorage) Create(ctx context.Context, text string) error {
	t := time.Now()
	result, err := ns.db.Exec(`INSERT INTO note (text, created_at, updated_at) VALUES (?, ?, ?)`, text, t, t)

	if err != nil {
		return err
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New("afftected 0 row")
	}

	return nil
}

func (ns *noteStorage) Update(ctx context.Context, id, text string) error {
	result, err := ns.db.Exec(`UPDATE note SET text=?, updated_at=? WHERE id=?`, text, time.Now(), id)

	if err != nil {
		return err
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New("afftected 0 row")
	}

	return nil
}

func (ns *noteStorage) Delete(ctx context.Context, id string) error {
	result, err := ns.db.Exec(`DELETE FROM note WHERE id=?`, id)

	if err != nil {
		return err
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New("afftected 0 row")
	}

	return nil
}

func (ns *noteStorage) Get(ctx context.Context, id string) (*models.Note, error) {
	note := &models.Note{}

	err := ns.db.QueryRow("SELECT id, text, created_at, updated_at FROM note WHERE id=?", id).Scan(&note.ID, &note.Text, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, errors.New("can't scan row")
	}

	return note, nil
}

func (ns *noteStorage) GetAll(ctx context.Context, orderBy string) ([]*models.Note, error) {
	if orderBy == "" {
		orderBy = "id"
	}

	notes := []*models.Note{}
	query := fmt.Sprintf("SELECT id, text, created_at, updated_at FROM note ORDER BY %s", orderBy)
	rows, err := ns.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.Text, &note.CreatedAt, &note.UpdatedAt)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}
	defer rows.Close()

	return notes, nil
}
