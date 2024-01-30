package v1

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"note/internal/models"
	"note/internal/models/mocks"
	"note/internal/service"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type tester struct {
	returning    *gomock.Call
	code         int
	errorMessage string
	reader       io.Reader
	isBadWriter  bool
	router       router
}

type router struct {
	method       string
	path         string
	body         string
	isHeaderNeed bool
	vars         map[string]string
}

type BadReader struct {
}

func (br BadReader) Read(p []byte) (int, error) {
	return 0, errors.New("some error")
}

type BadResponseWriter struct {
	http.ResponseWriter
}

func (bw BadResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("some error")
}

func getRequestRecorder(method, path, body string, vars map[string]string, reader io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()

	if reader != nil {
		req = httptest.NewRequest(method, path, reader)
	}

	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}

	return w, req
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vars := map[string]string{"id": "1"}

	srv := mocks.NewMockService(ctrl)
	handler := NewNoteHandler(srv, zap.L())
	req := httptest.NewRequest("*", "/", nil)
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()

	tes := []tester{
		{
			returning:    srv.EXPECT().Get(ctx, service.GetNote{ID: "1"}).Return(&models.Note{}, nil),
			code:         http.StatusOK,
			errorMessage: "expected 200, got:",
			router: router{
				method: "GET",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			returning:    srv.EXPECT().Get(ctx, service.GetNote{ID: "1"}).Return(nil, errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method: "GET",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "GET",
				path:   "/note/{id}",
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			isBadWriter:  true,
			router: router{
				method: "GET",
				path:   "/note/{id}",
			},
		},
		{
			returning:    srv.EXPECT().Get(ctx, service.GetNote{ID: "1"}).Return(nil, errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "GET",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			returning:    srv.EXPECT().Get(ctx, service.GetNote{ID: "1"}).Return(&models.Note{}, nil),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "GET",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
	}

	for _, test := range tes {
		w, req := getRequestRecorder(test.router.method, test.router.path, test.router.body, test.router.vars, test.reader)

		if test.isBadWriter {
			handler.GetByID(&BadResponseWriter{w}, req)
		} else {
			handler.GetByID(w, req)
		}

		if w.Code != test.code {
			t.Errorf(test.errorMessage+"%d", w.Code)
			return
		}
	}
}

func TestDeleteByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vars := map[string]string{"id": "1"}

	srv := mocks.NewMockService(ctrl)
	handler := NewNoteHandler(srv, zap.L())
	req := httptest.NewRequest("*", "/", nil)
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()

	tes := []tester{
		{
			returning:    srv.EXPECT().Delete(ctx, service.DeleteNote{ID: "1"}).Return(nil),
			code:         http.StatusOK,
			errorMessage: "expected 200, got:",
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			returning:    srv.EXPECT().Delete(ctx, service.DeleteNote{ID: "1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
				vars:   nil,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			isBadWriter:  true,
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
			},
		},
		{
			returning:    srv.EXPECT().Delete(ctx, service.DeleteNote{ID: "1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
		{
			returning:    srv.EXPECT().Delete(ctx, service.DeleteNote{ID: "1"}).Return(nil),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "DELETE",
				path:   "/note/{id}",
				vars:   vars,
			},
		},
	}

	for _, test := range tes {
		w, req := getRequestRecorder(test.router.method, test.router.path, test.router.body, test.router.vars, test.reader)

		if test.isBadWriter {
			handler.DeleteByID(&BadResponseWriter{w}, req)
		} else {
			handler.DeleteByID(w, req)
		}

		if w.Code != test.code {
			t.Errorf(test.errorMessage+"%d", w.Code)
			return
		}
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockService(ctrl)
	handler := NewNoteHandler(srv, zap.L())
	ctx := context.Background()

	tes := []tester{
		{
			returning:    srv.EXPECT().Create(ctx, service.CreateNote{Text: "test1"}).Return(nil),
			code:         http.StatusOK,
			errorMessage: "expected 200, got:",
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
		},
		{
			returning:    srv.EXPECT().Create(ctx, service.CreateNote{Text: "test1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
		},
		{
			returning:    srv.EXPECT().Create(ctx, service.CreateNote{Text: "test1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text" test1"}`,
				isHeaderNeed: true,
			},
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text" test1"}`,
				isHeaderNeed: true,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "POST",
				path:   "/note",
				body:   `{"text": "test1"}`,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "POST",
				path:   "/note",
				body:   `{"text": "test1"}`,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			reader:       &BadReader{},
			router: router{
				method:       "POST",
				path:         "/note",
				isHeaderNeed: true,
			},
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			reader:       &BadReader{},
			router: router{
				method:       "POST",
				path:         "/note",
				isHeaderNeed: true,
			},
			isBadWriter: true,
		},
		{
			returning:    srv.EXPECT().Create(ctx, service.CreateNote{Text: "test1"}).Return(nil),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method:       "POST",
				path:         "/note",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
		},
	}

	for _, test := range tes {
		w, req := getRequestRecorder(test.router.method, test.router.path, test.router.body, test.router.vars, test.reader)
		if test.router.isHeaderNeed {
			req.Header.Add("Content-type", "application/json")
		}

		if test.isBadWriter {
			handler.Create(&BadResponseWriter{w}, req)
		} else {
			handler.Create(w, req)
		}

		if w.Code != test.code {
			t.Errorf(test.errorMessage+"%d", w.Code)
			return
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vars := map[string]string{"id": "1"}

	srv := mocks.NewMockService(ctrl)
	handler := NewNoteHandler(srv, zap.L())
	req := httptest.NewRequest("*", "/", nil)
	req = mux.SetURLVars(req, vars)
	ctx := req.Context()

	tes := []tester{
		{
			returning:    srv.EXPECT().Update(ctx, service.UpdateNote{ID: "1", Text: "test1"}).Return(nil),
			code:         http.StatusOK,
			errorMessage: "expected 200, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				vars:         vars,
				isHeaderNeed: true,
			},
		},
		{
			returning:    srv.EXPECT().Update(ctx, service.UpdateNote{ID: "1", Text: "test1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				vars:         vars,
				isHeaderNeed: true,
			},
		},
		{
			returning:    srv.EXPECT().Update(ctx, service.UpdateNote{ID: "1", Text: "test1"}).Return(errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				vars:         vars,
				isHeaderNeed: true,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "PUT",
				path:   "/note/{id}",
				body:   `{"text": "test1"}`,
			},
		},
		{
			code:         http.StatusBadRequest,
			errorMessage: "expected 400, got:",
			router: router{
				method: "PUT",
				path:   "/note/{id}",
				body:   `{"text": "test1"}`,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text" test1"}`,
				isHeaderNeed: true,
				vars:         vars,
			},
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text" test1"}`,
				isHeaderNeed: true,
				vars:         vars,
			},
			isBadWriter: true,
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			reader:       &BadReader{},
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				isHeaderNeed: true,
				vars:         vars,
			},
		},
		{
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			reader:       &BadReader{},
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				isHeaderNeed: true,
				vars:         vars,
			},
			isBadWriter: true,
		},
		{
			returning:    srv.EXPECT().Update(ctx, service.UpdateNote{ID: "1", Text: "test1"}).Return(nil),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method:       "PUT",
				path:         "/note/{id}",
				body:         `{"text": "test1"}`,
				isHeaderNeed: true,
				vars:         vars,
			},
		},
	}

	for _, test := range tes {
		w, req := getRequestRecorder(test.router.method, test.router.path, test.router.body, test.router.vars, test.reader)
		if test.router.isHeaderNeed {
			req.Header.Add("Content-type", "application/json")
		}

		if test.isBadWriter {
			handler.UpdateByID(&BadResponseWriter{w}, req)
		} else {
			handler.UpdateByID(w, req)
		}

		if w.Code != test.code {
			t.Errorf(test.errorMessage+"%d", w.Code)
			return
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockService(ctrl)
	handler := NewNoteHandler(srv, zap.L())
	ctx := context.Background()

	tes := []tester{
		{
			returning:    srv.EXPECT().GetAll(ctx, service.GetNotes{OrderBy: "id"}).Return([]*models.Note{}, nil),
			code:         http.StatusOK,
			errorMessage: "expected 200, got:",
			router: router{
				method: "GET",
				path:   "/note?order_by=id",
			},
		},
		{
			returning:    srv.EXPECT().GetAll(ctx, service.GetNotes{OrderBy: "id"}).Return(nil, errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			router: router{
				method: "GET",
				path:   "/note?order_by=id",
			},
		},
		{
			returning:    srv.EXPECT().GetAll(ctx, service.GetNotes{OrderBy: "id"}).Return(nil, errors.New("some error")),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "GET",
				path:   "/note?order_by=id",
			},
		},
		{
			returning:    srv.EXPECT().GetAll(ctx, service.GetNotes{OrderBy: "id"}).Return([]*models.Note{}, nil),
			code:         http.StatusInternalServerError,
			errorMessage: "expected 500, got:",
			isBadWriter:  true,
			router: router{
				method: "GET",
				path:   "/note?order_by=id",
			},
		},
	}

	for _, test := range tes {
		w, req := getRequestRecorder(test.router.method, test.router.path, test.router.body, test.router.vars, test.reader)
		if test.router.isHeaderNeed {
			req.Header.Add("Content-type", "application/json")
		}

		if test.isBadWriter {
			handler.GetAll(&BadResponseWriter{w}, req)
		} else {
			handler.GetAll(w, req)
		}

		if w.Code != test.code {
			t.Errorf(test.errorMessage+"%d", w.Code)
			return
		}
	}
}
