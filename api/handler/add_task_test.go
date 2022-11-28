package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todotree/entity"
	"todotree/testutil"

	"github.com/go-playground/validator/v10"
)

func TestAddTask(t *testing.T) {
	t.Parallel()

	type wanted struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    wanted
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: wanted{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: wanted{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(
				ctx context.Context, title string,
			) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
