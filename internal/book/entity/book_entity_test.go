package entity_test

import (
	"testing"

	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	b, err := entity.New("American Gods", "Neil Gaiman", 100, 1)
	assert.Nil(t, err)
	assert.Equal(t, b.Title, "American Gods")
	assert.NotNil(t, b.ID)
}

func TestBook_Validate(t *testing.T) {
	type test struct {
		title    string
		author   string
		pages    int
		quantity int
		want     error
	}

	tests := []test{
		{
			title:    "American Gods",
			author:   "Neil Gaiman",
			pages:    100,
			quantity: 1,
			want:     nil,
		},
		{
			title:    "American Gods",
			author:   "Neil Gaiman",
			pages:    100,
			quantity: 0,
			want:     entity.ErrInvalidBookEntity,
		},
		{
			title:    "",
			author:   "Neil Gaiman",
			pages:    100,
			quantity: 1,
			want:     entity.ErrInvalidBookEntity,
		},
		{
			title:    "American Gods",
			author:   "",
			pages:    100,
			quantity: 1,
			want:     entity.ErrInvalidBookEntity,
		},
		{
			title:    "American Gods",
			author:   "Neil Gaiman",
			pages:    0,
			quantity: 1,
			want:     entity.ErrInvalidBookEntity,
		},
	}
	for _, tc := range tests {

		_, err := entity.New(tc.title, tc.author, tc.pages, tc.quantity)
		assert.Equal(t, err, tc.want)
	}

}
