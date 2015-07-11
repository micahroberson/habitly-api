package models

import (
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Habit struct {
  Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name        string        `json:"name" bson:"name"`
  AccountId   bson.ObjectId `json:"account_id" bson:"account_id"`
  CreatedAt   int64         `json:"created_at" bson:"created_at"`
  UpdatedAt   int64         `json:"updated_at" bson:"updated_at"`
}

type HabitResource struct {
  Payload Habit `json:"payload"`
}

type HabitCollection struct {
  Payload []Habit `json:"payload"`
}

type HabitRepo struct {
  Coll *mgo.Collection
}

func (r *HabitRepo) All(q bson.M) (HabitCollection, error) {
  result := HabitCollection{[]Habit{}}
  err := r.Coll.Find(q).Sort("created_at").All(&result.Payload)
  if err != nil {
    return result, err
  }
  return result, nil
}

func (r *HabitRepo) Find(q bson.M) (HabitResource, error) {
  result := HabitResource{}
  err := r.Coll.Find(q).One(&result.Payload)
  if err != nil {
    return result, err
  }
  return result, nil
}

func (r *HabitRepo) Create(habit *Habit) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, habit)
  if err != nil {
    return err
  }
  habit.Id = id

  return nil
}

func (r *HabitRepo) Update(q bson.M, update bson.M) error {
  err := r.Coll.Update(q, update)
  if err != nil {
    return err
  }

  return nil
}

func (r *HabitRepo) Delete(q bson.M) error {
  err := r.Coll.Remove(q)
  if err != nil {
    return err
  }

  return nil
}