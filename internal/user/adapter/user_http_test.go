package adapter_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sgraham785/gocleanarch-example/internal/user/adapter"
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/mock"
	"github.com/stretchr/testify/assert"
)

func Test_ListUsersHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.UserRouter(r, *n, m)
	path, err := r.GetRoute("listUsers").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user", path)
	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().
		ListUsers().
		Return([]*entity.User{u}, nil)
	ts := httptest.NewServer(adapter.ListUsersHTTP(m))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_ListUsersHTTP_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	ts := httptest.NewServer(adapter.ListUsersHTTP(m))
	defer ts.Close()
	m.EXPECT().
		SearchUsers("dio").
		Return(nil, entity.ErrUserNotFound)
	res, err := http.Get(ts.URL + "?name=dio")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_ListUsersHTTP_Search(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().
		SearchUsers("ozzy").
		Return([]*entity.User{u}, nil)
	ts := httptest.NewServer(adapter.ListUsersHTTP(m))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?name=ozzy")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_CreateUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.UserRouter(r, *n, m)
	path, err := r.GetRoute("createUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user", path)

	m.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.NewID(), nil)
	h := adapter.CreateUserHTTP(m)

	ts := httptest.NewServer(h)
	defer ts.Close()
	payload := fmt.Sprintf(`{
		"name": "ozzy",
		"email": "ozzy@hell.com",
		"password": "asasa",
		"first_name":"Ozzy",
		"last_name":"Osbourne"
		}`)
	resp, _ := http.Post(ts.URL+"/v1/user", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var u *adapter.UserHTTP
	json.NewDecoder(resp.Body).Decode(&u)
	assert.Equal(t, "Ozzy Osbourne", fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

func Test_GetUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.UserRouter(r, *n, m)
	path, err := r.GetRoute("getUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user/{id}", path)
	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().
		GetUser(u.ID).
		Return(u, nil)
	handler := adapter.GetUserHTTP(m)
	r.Handle("/v1/user/{id}", handler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, err := http.Get(ts.URL + "/v1/user/" + u.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *adapter.UserHTTP
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, u.ID, d.ID)
}

func Test_DeleteUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.UserRouter(r, *n, m)
	path, err := r.GetRoute("deleteUser").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/user/{id}", path)
	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().DeleteUser(u.ID).Return(nil)
	handler := adapter.DeleteUserHTTP(m)
	req, _ := http.NewRequest("DELETE", "/v1/user/"+u.ID.String(), nil)
	r.Handle("/v1/user/{id}", handler).Methods("DELETE", "OPTIONS")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
