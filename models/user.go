package models

import (
  // "github.com/micahroberson/habitly-api/lib"
  "github.com/dchest/uniuri"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
)

type User struct {
  Id               bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name             string        `json:"name" bson:"name"`
  Password         []byte        `json:"password" bson:"password"`
  Email            string        `json:"email" bson:"email"`
  CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
  LastLoggedInAt   time.Time     `json:"last_logged_in_at" bson:"last_logged_in_at"`
  AuthToken        string        `json:"auth_token" bson:"auth_token"`
}

// type UserLogin struct {
//   Email            string        `json:"email" binding:"required"`
//   Password         string        `json:"password" binding:"required"`
// }

func (u *User) Login(db *mgo.Database) {
  u.LastLoggedInAt = time.Now()
  u.AuthToken = uniuri.NewLen(22)

  err:= db.C("users").Update(bson.M{"_id": u.Id}, u)
  if err != nil {
    panic(err)
  }
  // u.authenticated = true
}

func (u *User) Logout() {
  // u.authenticated = false
}

// func NewAuth() Authenticator {
//   return &User{Email: ""}
// }