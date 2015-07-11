package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type User struct {
  Id               bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name             string        `json:"name" bson:"name"`
  Password         []byte        `json:"password" bson:"password"`
  Email            string        `json:"email" bson:"email"`
  CreatedAt        int64         `json:"created_at" bson:"created_at"`
  LastSignedInAt   int64         `json:"last_signed_in_at" bson:"last_signed_in_at"`
}

type UserResource struct {
  Payload User `json:"payload"`
}

type UserCollection struct {
  Payload []User `json:"payload"`
}

type UserRepo struct {
  Coll *mgo.Collection
}

func (r *UserRepo) Find(q bson.M) (UserResource, error) {
  result := UserResource{}
  err := r.Coll.Find(q).One(&result.Payload)
  if err != nil {
    return result, err
  }
  return result, nil
}

func (r *UserRepo) Create(user *User) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, user)
  if err != nil {
    return err
  }
  user.Id = id

  return nil
}

func (r *UserRepo) Update(q bson.M, update bson.M) error {
  err := r.Coll.Update(q, update)
  if err != nil {
    return err
  }

  return nil
}