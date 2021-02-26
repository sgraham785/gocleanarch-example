package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/xid"
	"github.com/sgraham785/gocleanarch-example/internal/user/usecase"

	"github.com/sgraham785/gocleanarch-example/internal/user/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// UserHTTP JSON data
type UserHTTP struct {
	ID        entity.ID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// ListUsersHTTP handler
func ListUsersHTTP(u usecase.UserUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading users"
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
func CreateUserHTTP(u usecase.UserUseCase) http.Handler {
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
func GetUserHTTP(u usecase.UserUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading user"
		vars := mux.Vars(r)
		id, err := xid.FromString(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := u.GetUser(id)
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
	})
}

// DeleteUserHTTP handler
func DeleteUserHTTP(u usecase.UserUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing user"
		vars := mux.Vars(r)
		id, err := xid.FromString(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = u.DeleteUser(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// UserRouter defined http routes for handlers
func UserRouter(r *mux.Router, n negroni.Negroni, u usecase.UserUseCase) {
	r.Handle("/v1/user", n.With(
		negroni.Wrap(ListUsersHTTP(u)),
	)).Methods("GET", "OPTIONS").Name("listUsers")

	r.Handle("/v1/user", n.With(
		negroni.Wrap(CreateUserHTTP(u)),
	)).Methods("POST", "OPTIONS").Name("createUser")

	r.Handle("/v1/user/{id}", n.With(
		negroni.Wrap(GetUserHTTP(u)),
	)).Methods("GET", "OPTIONS").Name("getUser")

	r.Handle("/v1/user/{id}", n.With(
		negroni.Wrap(DeleteUserHTTP(u)),
	)).Methods("DELETE", "OPTIONS").Name("deleteUser")
}
