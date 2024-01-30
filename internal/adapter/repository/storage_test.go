package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"note/internal/models"
	"reflect"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("can't create mock: %s", err)
		return
	}
	defer db.Close()

	var id string = "1"
	var ctx context.Context

	ti := time.Now()

	rows := getRows(1, "text message", ti)
	expect := &models.Note{
		ID:        1,
		Text:      "text message",
		CreatedAt: ti,
		UpdatedAt: ti,
	}
	repo := NewStorage(db)

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note WHERE").WithArgs(id).WillReturnRows(rows)

	note, err := repo.Get(ctx, id)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if !reflect.DeepEqual(note, expect) {
		t.Errorf("results not match, want %v, have %v", expect, note)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note WHERE").WithArgs(id).WillReturnError(errors.New("some error"))

	_, err = repo.Get(ctx, "1")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	_, err = repo.Get(ctx, "")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}
}

func testCUDperation(t *testing.T, operationName string, operation func() error, mock sqlmock.Sqlmock) {
	mock.ExpectExec(operationName).WillReturnResult(sqlmock.NewResult(1, 1))

	err := operation()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectExec(operationName).WillReturnError(errors.New("some error"))

	err = operation()

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectExec(operationName).WillReturnResult(sqlmock.NewResult(0, 0))

	err = operation()

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("can't create mock: %s", err)
		return
	}
	defer db.Close()

	var ctx context.Context
	repo := NewStorage(db)

	operation := func() error {
		return repo.Create(ctx, "message")
	}
	testCUDperation(t, "INSERT INTO note", operation, mock)

	err = repo.Create(ctx, "")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("can't create mock: %s", err)
		return
	}
	defer db.Close()

	var ctx context.Context
	repo := NewStorage(db)

	operation := func() error {
		return repo.Update(ctx, "1", "message")
	}
	testCUDperation(t, "UPDATE note", operation, mock)

	err = repo.Update(ctx, "1", "")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	err = repo.Update(ctx, "", "message")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("can't create mock: %s", err)
		return
	}
	defer db.Close()

	var ctx context.Context
	repo := NewStorage(db)

	operation := func() error {
		return repo.Delete(ctx, "1")
	}
	testCUDperation(t, "DELETE FROM note WHERE", operation, mock)

	err = repo.Delete(ctx, "")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}
}

func getRows(count int, text string, ti time.Time) *sqlmock.Rows {

	rows := sqlmock.NewRows([]string{"id", "text", "created_at", "updated_at"})
	for i := 0; i < count; i++ {
		rows.AddRow(i+1, text, ti, ti)
	}

	return rows
}

func TestGetList(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("can't create mock: %s", err)
		return
	}
	defer db.Close()

	var ctx context.Context
	repo := NewStorage(db)
	ti := time.Now()

	rows := getRows(3, "text message", ti)
	expect := []*models.Note{
		{
			ID:        1,
			Text:      "text message",
			CreatedAt: ti,
			UpdatedAt: ti,
		},
		{
			ID:        2,
			Text:      "text message",
			CreatedAt: ti,
			UpdatedAt: ti,
		},
		{
			ID:        3,
			Text:      "text message",
			CreatedAt: ti,
			UpdatedAt: ti,
		},
	}

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note ORDER BY").WillReturnRows(rows)

	notes, err := repo.GetAll(ctx, "id")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if !reflect.DeepEqual(expect, notes) {
		t.Errorf("results not match, want %v, have %v", expect, notes)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note ORDER BY").WillReturnError(errors.New("some error"))

	_, err = repo.GetAll(ctx, "id")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	errorRows := sqlmock.NewRows([]string{"time"}).AddRow(time.Now()).AddRow(time.Now())

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note ORDER BY").WillReturnRows(errorRows)

	_, err = repo.GetAll(ctx, "id")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	rows = getRows(3, "text message", ti)

	mock.ExpectQuery("SELECT id, text, created_at, updated_at FROM note ORDER BY").WillReturnRows(rows)

	notes, err = repo.GetAll(ctx, "")

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if !reflect.DeepEqual(expect, notes) {
		t.Errorf("results not match, want %v, have %v", expect, notes)
		return
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	_, err = repo.GetAll(ctx, "test")

	if err == nil {
		t.Error("expected error, got nil")
		return
	}
}
