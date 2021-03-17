package adapter_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sgraham785/gocleanarch-example/internal/book/adapter"
	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/mock"
	"github.com/sgraham785/gocleanarch-example/pkg/router"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

func TestListBooksHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	u := mock.NewMockBookUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	b := &entity.Book{
		ID: entity.NewID(),
	}
	u.EXPECT().
		ListBooks().
		Return([]*entity.Book{b}, nil)

	adapter.HTTPRoutes(s, u)
	h := adapter.ListBooksHTTP(u)
	ts := httptest.NewServer(h)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListBooksHTTP_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	u := mock.NewMockBookUseCase(controller)
	h := adapter.ListBooksHTTP(u)

	ts := httptest.NewServer(h)
	defer ts.Close()

	u.EXPECT().
		SearchBooks("book of books").
		Return(nil, entity.ErrBookNotFound)

	res, err := http.Get(ts.URL + "?title=book+of+books")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// func TestListBooks_Search(t *testing.T) {
// 	controller := gomock.NewController(t)
// 	defer controller.Finish()
// 	u := mock.NewMockBookUseCase(controller)
// 	b := &entity.Book{
// 		ID: entity.NewID(),
// 	}
// 	u.EXPECT().
// 		SearchBooks("ozzy").
// 		Return([]*entity.Book{b}, nil)
// 	ts := httptest.NewServer(adapter.ListBooksHTTP(u))
// 	defer ts.Close()
// 	res, err := http.Get(ts.URL + "?title=ozzy")
// 	assert.Nil(t, err)
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// }

func TestCreateBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	u := mock.NewMockBookUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	u.EXPECT().
		CreateBook(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.NewID(), nil)

	adapter.HTTPRoutes(s, u)
	h := adapter.CreateBookHTTP(u)
	ts := httptest.NewServer(h)
	defer ts.Close()
	payload := fmt.Sprintf(`{
		"title": "I Am Ozzy",
		"author": "Ozzy Osbourne",
		"pages": 294,
		"quantity":1
	}`)

	resp, _ := http.Post(ts.URL+"/book", "application/json", strings.NewReader(payload))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var b *entity.Book
	json.NewDecoder(resp.Body).Decode(&b)
	assert.Equal(t, "Ozzy Osbourne", b.Author)
}

func TestGetBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	u := mock.NewMockBookUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	b := &entity.Book{
		ID: entity.NewID(),
	}

	u.EXPECT().
		GetBook(b.ID.String()).
		Return(b, nil)

	adapter.HTTPRoutes(s, u)
	h := adapter.GetBookHTTP(u)
	r.Chi.Handle("/book/{bookID}", h)
	ts := httptest.NewServer(r.Chi)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/book/" + b.ID.String() + ".json")

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var d *entity.Book
	json.NewDecoder(res.Body).Decode(&d)
	assert.NotNil(t, d)
	assert.Equal(t, b.ID, d.ID)
}

func TestDeleteBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	u := mock.NewMockBookUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}

	b := &entity.Book{
		ID: entity.NewID(),
	}
	u.EXPECT().
		DeleteBook(b.ID.String()).
		Return(nil)

	adapter.HTTPRoutes(s, u)
	h := adapter.DeleteBookHTTP(u)
	r.Chi.Handle("/book/{bookID}", h)
	ts := httptest.NewServer(r.Chi)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", "/book/"+b.ID.String(), nil)
	rr := httptest.NewRecorder()
	r.Chi.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
