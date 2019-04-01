package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"topos-backend-assignment/db"
)

type UserLogin struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Password string        `bson:"password" json:"password"`
	Username string        `bson:"username" json:"username"`
}

func getUserLoginCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C(db.UserLoginCollection)
}

/**
Helper method to create a new user
*/
func (user *UserLogin) CreateUser(session *mgo.Session, u UserLogin) error {
	err := getUserLoginCollection(session).Insert(&u)
	return err
}

/**
Helper method for user login verification
*/
func (user *UserLogin) LoginUser(session *mgo.Session, uname string) (string, error) {
	var u UserLogin
	err := getUserLoginCollection(session).Find(bson.M{"username": uname}).One(&u)
	return u.Password, err
}
