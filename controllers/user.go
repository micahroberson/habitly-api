package controllers

import (
  "github.com/micahroberson/habitly-api/models"
  // "github.com/micahroberson/habitly-api/lib"
  "code.google.com/p/go.crypto/bcrypt"
  "github.com/dchest/uniuri"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/go-martini/martini"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  // "fmt"
)

func RegisterUser(db *mgo.Database, r render.Render, user models.User) {
  user.Id = bson.NewObjectId()

  user.CreatedAt = time.Now().UTC()
  bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
  user.Password = bcryptPassword
  user.LastLoggedInAt = time.Now()
  user.AuthToken = uniuri.NewLen(22)
  // salt := lib.GenerateSalt()
  // hash := lib.EncodePbkdf2(user.Password, salt)
  // user.Password = salt + hash

  err := db.C("users").Insert(user)
  if err != nil {
    panic(err)
  }

  r.JSON(201, user)
}

func LoginUser(db *mgo.Database, r render.Render, u models.User) {
  var user models.User
  
  err := db.C("users").Find(bson.M{"email": u.Email}).One(&user)

  if err != nil {
    r.JSON(401, nil)
  } else {
    passwordErr := bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))
    if passwordErr == nil {
      r.JSON(200, user)
    } else {
      r.JSON(401, nil)
    }
  }
}

func GetUser(db *mgo.Database, r render.Render, p martini.Params) {
  var user models.User

  err := db.C("users").Find(bson.M{"email": p["email"]}).One(&user)

  // TODO Check for 404

  if err != nil {
    panic(err)
  }

  r.JSON(200, user)
}