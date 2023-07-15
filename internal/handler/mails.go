package handler

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/internal/validation"
	"github.com/amaretur/mail-client/pkg/mails"
)

type MailUsecase interface {
	GetFolders(uid uint) ([]mails.Mailbox, error)
	GetMails(uid uint, folder string,
		offset, limit uint) ([]*mails.Mail, uint32, error)
	SendMail(uid uint, encrypt, sign bool, mail *mails.Mail) error
}

type MailHandler struct {
	usecase	MailUsecase
	amw		AuthMiddleware
}

func NewMailHandler(usecase MailUsecase, amw AuthMiddleware) *MailHandler {
	return &MailHandler{
		usecase: usecase,
		amw: amw,
	}
}

func (m *MailHandler) Init(router *mux.Router) {

	authonly := router.PathPrefix("").Subrouter()
	authonly.Use(m.amw.AuthOnly)
	authonly.HandleFunc("/folder", m.getFolders).Methods("GET")

	authonly.HandleFunc("", m.getMails).
		Methods("GET").Queries("folder", "{folder}")
	authonly.HandleFunc("", m.sendMail).
		Methods("POST")

}

func (m *MailHandler) getFolders(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	folders, err := m.usecase.GetFolders(uid)
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	NewJSONResponse(w, folders)
}

func (m *MailHandler) getMails(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	offset, limit := getOffsetAndLimit(r)
	if limit == 0 {
		limit = 10 // default limit
	}

	folderId := r.URL.Query().Get("folder")

	mails, total, err := m.usecase.GetMails(
		uid, folderId, uint(offset), uint(limit))
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	NewJSONResponse(w, &dto.Mails{Mails: mails, Total: total})
}

func (m *MailHandler) sendMail(w http.ResponseWriter, r* http.Request) {

	uid := r.Context().Value("uid").(uint)

	var data struct {
		Encrypt	bool		`json:"encrypt" binding:"required"`
		Sign	bool		`json:"sign" binding:"required"`
		Mail	mails.Mail	`json:"mail"`
	}

	err := json.NewDecoder(r.Body).Decode(&data);
	if  err != nil || !validation.IsValidMsg(&data.Mail) {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	err = m.usecase.SendMail(uid, data.Encrypt, data.Sign, &data.Mail)
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
