package user

import (
	"net/http"

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
	router.GET(usersURL, h.GetList)
	router.POST(usersURL, h.CreateUser)
	router.GET(userURL, h.GetUserByUIID)
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartiallyUpdateUser)
	router.DELETE(userURL, h.DeleteUser)

}

func (h handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is list of users"))
	w.WriteHeader(200)
}

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is create of users"))
	w.WriteHeader(204)
}

func (h handler) GetUserByUIID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is get user by UIID of users"))
	w.WriteHeader(200)
}

func (h handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is update  users"))
	w.WriteHeader(200)
}

func (h handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is partially update  users"))
	w.WriteHeader(200)
}

func (h handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("This is delete users"))
	w.WriteHeader(204)
}
