package models

type (
	User struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
	}
)
