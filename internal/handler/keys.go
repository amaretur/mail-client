package handler

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/internal/validation"
)

type KeysUsecase interface {
	GetKeysList(uid, offset, limit uint) (*dto.DialogKeysList, error)
	GetKeys(id uint) (*dto.DialogKeys, error)
	GetVerify(uid uint) (*dto.VerifyKeys, error)
	UpdateVerify(uid uint, keys *dto.VerifyKeys) error
	SetRandomVerify(uid uint) error
	Receive(uid uint, data *dto.Receive) (uint, error)
	Share(uid uint, interlocutor string) error
}

type KeyHandler struct {
	usecase	KeysUsecase
	amw		AuthMiddleware
}

func NewKeyHandler(usecase KeysUsecase, amw AuthMiddleware) *KeyHandler {
	return &KeyHandler{
		usecase: usecase,
		amw: amw,
	}
}

func (k *KeyHandler) Init(router *mux.Router) {

	authonly := router.PathPrefix("").Subrouter()
	authonly.Use(k.amw.AuthOnly)

	authonly.HandleFunc("", k.getKeysList).Methods("GET")
	authonly.HandleFunc("/{id:[0-9]+}", k.getKeys).Methods("GET")
	authonly.HandleFunc("/verify", k.getVerify).Methods("GET")
	authonly.HandleFunc("/verify", k.setRandomVerify).
		Methods("PUT").
		Queries("random", "")
	authonly.HandleFunc("/verify", k.updateVerify).Methods("PUT")

	authonly.HandleFunc("/share", k.share).Methods("POST")
	authonly.HandleFunc("/receive", k.receive).Methods("POST")
}

func (k *KeyHandler) getKeysList(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	offset, limit := getOffsetAndLimit(r)
	if limit == 0 {
		limit = 10 // default limit
	}

	data, err := k.usecase.GetKeysList(uid, uint(offset), uint(limit))
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	NewJSONResponse(w, data)
}

func (k *KeyHandler) getKeys(w http.ResponseWriter, r* http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	keys, err := k.usecase.GetKeys(uint(id))
	if err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	NewJSONResponse(w, keys)
}

func (k *KeyHandler) getVerify(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	verify, err := k.usecase.GetVerify(uid)
	if err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	NewJSONResponse(w, verify)
}

func (k *KeyHandler) updateVerify(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	data := new(dto.VerifyKeys)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	if err := k.usecase.UpdateVerify(uid, data); err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (k *KeyHandler) setRandomVerify(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	if err := k.usecase.SetRandomVerify(uid); err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (k *KeyHandler) share(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	var data struct {
		Interlocutor string	`json:"recipient" binding:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if  err != nil || !validation.IsValidMail(data.Interlocutor) {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	if err := k.usecase.Share(uid, data.Interlocutor); err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (k *KeyHandler) receive(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	data := new(dto.Receive)

	err := json.NewDecoder(r.Body).Decode(data);
	if err != nil || !validation.IsValidReceive(data) {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	id, err := k.usecase.Receive(uid, data)
	if err != nil {
		NewJSONResponseWhithError(w, http.StatusNotFound, "not found")
		return
	}

	NewJSONResponse(w, map[string]uint{"id":id})
}


