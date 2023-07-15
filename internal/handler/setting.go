package handler

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/internal/validation"
)

type SettingUsecase interface {
	GetSetting(uid uint) (*dto.Setting, error)
	UpdateSetting(data *dto.Setting) error
}

type SettingHandler struct {
	usecase	SettingUsecase
	amw		AuthMiddleware
}

func NewSettingHandler(usecase SettingUsecase,
	amw AuthMiddleware) *SettingHandler {

	return &SettingHandler{
		usecase: usecase,
		amw: amw,
	}
}

func (s *SettingHandler) Init(router *mux.Router) {

	authonly := router.PathPrefix("").Subrouter()
	authonly.Use(s.amw.AuthOnly)
	authonly.HandleFunc("", s.get).Methods("GET")
	authonly.HandleFunc("", s.update).Methods("PUT")
}

func (s *SettingHandler) get(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	setting, err := s.usecase.GetSetting(uid)
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	NewJSONResponse(w, setting)
}

func (s *SettingHandler) update(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	setting := new(dto.Setting)

	err := json.NewDecoder(r.Body).Decode(setting); 
	if err != nil || !validation.IsValidSetting(setting) {
		NewJSONResponseWhithError(w, http.StatusBadRequest, "")
		return
	}

	setting.UserId = uid

	err = s.usecase.UpdateSetting(setting)
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
