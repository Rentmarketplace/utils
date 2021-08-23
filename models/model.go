package models

type Error struct {
	Message string `json:"message"`
	Code int16 `json:"code"`
}

type JWT struct {
	Authorization string `json:"Authorization" binding:"required"`
}
