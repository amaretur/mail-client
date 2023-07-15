package dto

type AuthData struct {
	Username	string	`json:"username" binding:"required"`
	Password	string	`json:"password" binding:"required"`
	Setting		Setting	`json:"setting" binding:"required"`
}


