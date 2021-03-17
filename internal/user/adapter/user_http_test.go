package adapter_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sgraham785/gocleanarch-example/internal/user/adapter"
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/mock"
	"github.com/sgraham785/gocleanarch-example/pkg/router"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestListUsersHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	m := mock.NewMockUserUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().
		ListUsers().
		Return([]*entity.User{u}, nil)

	adapter.HTTPRoutes(s, m)
	h := adapter.ListUsersHTTP(m)

	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListUsersHTTP_NotFound(t *testing.T) {
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

func TestListUsersHTTP_Search(t *testing.T) {
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

func TestCreateUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	m.EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.NewID(), nil)

	adapter.HTTPRoutes(s, m)
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
	resp, _ := http.Post(ts.URL+"/user", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var u *adapter.UserHTTP
	json.NewDecoder(resp.Body).Decode(&u)
	assert.Equal(t, "Ozzy Osbourne", fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

func TestGetUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	u := &entity.User{
		ID: entity.NewID(),
	}
	m.EXPECT().
		GetUser(u.ID.String()).
		Return(u, nil)

	adapter.HTTPRoutes(s, m)
	h := adapter.GetUserHTTP(m)
	r.Chi.Handle("/user/{userID}", h)
	ts := httptest.NewServer(r.Chi)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/user/" + u.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *adapter.UserHTTP
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, u.ID, d.ID)
}

func TestDeleteUserHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	m := mock.NewMockUserUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	u := &entity.User{
		ID: entity.NewID(),
	}

	m.EXPECT().
		DeleteUser(u.ID.String()).
		Return(nil)

	adapter.HTTPRoutes(s, m)
	h := adapter.DeleteUserHTTP(m)
	r.Chi.Handle("/user/{userID}", h)
	ts := httptest.NewServer(r.Chi)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", "/user/"+u.ID.String(), nil)
	rr := httptest.NewRecorder()
	r.Chi.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
