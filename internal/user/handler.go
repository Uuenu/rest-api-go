package user

import (
	"fmt"
	"net/http"

	"rest-api-go/internal/apperror"
	"rest-api-go/internal/handlers"
	"rest-api-go/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{} //

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUIID)) // ??????????????????????
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))

}

func (h handler) GetList(w http.ResponseWriter, r *http.Request) error {
	return apperror.ErrNotFound
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("this is API error")
}

func (h handler) GetUserByUIID(w http.ResponseWriter, r *http.Request) error {
	return apperror.NewAppError(nil, "test", "devel text", "TEST-000005")
}

func (h handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is update  users"))
	w.WriteHeader(200)

	return nil
}

func (h handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is partially update  users"))
	w.WriteHeader(200)

	return nil
}

func (h handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is delete users"))
	w.WriteHeader(204)

	return nil
}
