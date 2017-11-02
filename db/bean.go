package db

import (
	"gopkg.in/mgo.v2/bson"
)

type AccountBean struct {
	AccountId   bson.ObjectId `bson:"_id"`
	AccountName string        `bson:"name"`
	Password    string        `bson:"password"`
	AccountType int           `bson:"type"`
	CreateDate  string        `bson:"create_date"`
}

type UserBean struct {
	UserId      bson.ObjectId `bson:"_id"`
	UserName    string        `bson:"name"`
	DisplayName string        `bson:"display_name"`
	Phone       string        `bson:"phone"`
	Email       string        `bson:"email"`
	Comments    string        `bson:"comments"`
	Password    string        `bson:"password"`
	Account     string        `bson:"account"`
	CreateDate  string        `bson:"create_date"`
}

type GroupBean struct {
	GroupId    bson.ObjectId `bson:"_id"`
	GroupName  string        `bson:"name"`
	Comments   string        `bson:"comments"`
	Account    string        `bson:"account"`
	CreateDate string        `bson:"create_date"`
}

type KeyBean struct {
	KeyId       string `bson:"_id"`
	CreatorType int    `bson:"creator_type"`
}
