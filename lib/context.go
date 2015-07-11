package lib

import (
  "crypto/rsa"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type AppContext struct {
  MongoSession     *mgo.Database
  VerifyKey        *rsa.PublicKey
  SignKey          *rsa.PrivateKey
  currentUserId    bson.ObjectId
}

func (c *AppContext) SetCurrentUserId(id string) {
  c.currentUserId = bson.ObjectIdHex(id)
}

func (c AppContext) GetCurrentUserId() (bson.ObjectId) {
  return c.currentUserId
}