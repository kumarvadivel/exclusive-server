package models

type Todo struct{
	ID string `json:"_id",bson:"_id"`
	content string `json:"content",bson:"content"`
}