package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

// UserHTTP JSON data
type UserHTTP struct {
	ID        entity.ID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// ListUsersHTTP handler
func ListUsersHTTP(u usecase.UserUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error listing users"
		var data []*entity.User
		var err error
		name := r.URL.Query().Get("name")
		switch {
		case name == "":
			data, err = u.ListUsers()
		default:
			data, err = u.SearchUsers(name)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrUserNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*UserHTTP
		for _, d := range data {
			toJ = append(toJ, &UserHTTP{
				ID:        d.ID,
				Email:     d.Email,
				FirstName: d.FirstName,
				LastName:  d.LastName,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

// CreateUserHTTP handler
func CreateUserHTTP(u usecase.UserUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding user"
		var input struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		id, err := u.CreateUser(input.Email, input.Password, input.FirstName, input.LastName)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &UserHTTP{
			ID:        id,
			Email:     input.Email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// GetUserHTTP handler
func GetUserHTTP(u usecase.UserUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading user"
		if userID := chi.URLParam(r, "userID"); userID != "" {
			data, err := u.GetUser(userID)
			w.Header().Set("Content-Type", "application/json")
			if err != nil && err != entity.ErrUserNotFound {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error getting user"))
				return
			}

			if data == nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errorMessage))
				return
			}
			toJ := &UserHTTP{
				ID:        data.ID,
				Email:     data.Email,
				FirstName: data.FirstName,
				LastName:  data.LastName,
			}
			if err := json.NewEncoder(w).Encode(toJ); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// DeleteUserHTTP handler
func DeleteUserHTTP(u usecase.UserUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing user"
		if userID := chi.URLParam(r, "userID"); userID != "" {
			err := u.DeleteUser(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
		}
	})
}

// HTTPRoutes defined http routes for user
func HTTPRoutes(s *server.Server, u usecase.UserUseCase) {
	// RESTy routes for "books" resource
	s.Router.Chi.Route("/user", func(r chi.Router) {
		// r.With(paginate).Get("/", ListBooksHTTP(u))
		r.Get("/", ListUsersHTTP(u))
		r.Post("/", CreateUserHTTP(u))     // POST /book
		r.Get("/search", ListUsersHTTP(u)) // GET /book/search?title=something

		r.Route("/{userID}", func(r chi.Router) {
			// r.Use(BookCtx)             // Load the *Article on the request context
			r.Get("/", GetUserHTTP(u)) // GET /book/123
			// r.Put("/", UpdateArticle)    // PUT /book/123
			r.Delete("/", DeleteUserHTTP(u)) // DELETE /book/123
		})

		// GET /book/whats-up
		// r.With(BookCtx).Get("/{bookTitle:[a-z-]+}", GetBookHTTP(u))
	})
}
