package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "gopkg.in/mgo.v2/bson"
  "golang.org/x/crypto/bcrypt"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
  "github.com/micahroberson/habitly-api/lib"
  "github.com/micahroberson/habitly-api/models"
)

func AuthenticateUser(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("AuthenticateUser")

  body := context.Get(r, "body").(*models.UserResource)
  repo := models.UserRepo{c.MongoSession.C("users")}

  q := bson.M{
    "email": body.Payload.Email,
  }

  userResource, err := repo.Find(q)

  if err != nil {
    panic(err)
  }

  if err := bcrypt.CompareHashAndPassword(userResource.Payload.Password, []byte(body.Payload.Password)); err != nil {
    panic(err)
  }

  token := jwt.New(jwt.GetSigningMethod("RS256"))
  token.Claims["AccessToken"] = "level1"
  token.Claims["name"] = userResource.Payload.Name
  token.Claims["user_id"] = userResource.Payload.Id
  tokenString, err := token.SignedString(c.SignKey)
  if err != nil {
    panic(err)
  }

  update := bson.M{
    "$set": bson.M{
      "last_signed_in_at": lib.MakeTimestamp(),
    },
  }

  if err := repo.Update(q, update); err != nil {
    panic(err)
  }

  w.Header().Set("X-Access-Token", tokenString)

  if err := json.NewEncoder(w).Encode(userResource); err != nil {
    panic(err)
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("\n"))
}

func RegisterUser(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("RegisterUser")

  body := context.Get(r, "body").(*models.UserResource)
  repo := models.UserRepo{c.MongoSession.C("users")}

  body.Payload.CreatedAt = lib.MakeTimestamp()
  bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Payload.Password), bcrypt.DefaultCost)
  body.Payload.Password = bcryptPassword
  body.Payload.LastSignedInAt = lib.MakeTimestamp()

  if err := repo.Create(&body.Payload); err != nil {
    panic(err)
  }

  token := jwt.New(jwt.GetSigningMethod("RS256"))
  token.Claims["AccessToken"] = "level1"
  token.Claims["name"] = body.Payload.Name
  token.Claims["user_id"] = body.Payload.Id
  tokenString, err := token.SignedString(c.SignKey)

  if err != nil {
    panic(err)
  }

  w.Header().Set("X-Access-Token", tokenString)
  w.WriteHeader(201)
  json.NewEncoder(w).Encode(body)
}
