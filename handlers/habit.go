package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  "github.com/julienschmidt/httprouter"
  "github.com/micahroberson/habitly-api/lib"
  "github.com/micahroberson/habitly-api/models"
  // "time"
)

func GetAllHabits(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("GetAllHabits")

  q := bson.M{
    "account_id": c.GetCurrentUserId(),
  }

  repo := models.HabitRepo{c.MongoSession.C("habits")}
  habits, err := repo.All(q)

  if err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(habits)
}

func GetHabit(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("GetHabit")

  params := context.Get(r, "params").(httprouter.Params)

  q := bson.M{
    "account_id": c.GetCurrentUserId(),
    "_id":        bson.ObjectIdHex(params.ByName("id")),
  }

  repo := models.HabitRepo{c.MongoSession.C("habits")}
  habit, err := repo.Find(q)

  if err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(habit)
}

func CreateHabit(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("CreateHabit")

  body := context.Get(r, "body").(*models.HabitResource)
  repo := models.HabitRepo{c.MongoSession.C("habits")}

  // Set account_id
  body.Payload.AccountId = c.GetCurrentUserId()
  body.Payload.CreatedAt = lib.MakeTimestamp()
  body.Payload.UpdatedAt = lib.MakeTimestamp()

  err := repo.Create(&body.Payload)

  if err != nil {
    panic(err)
  }

  w.WriteHeader(201)
  json.NewEncoder(w).Encode(body)
}

func UpdateHabit(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("UpdateHabit")

  params := context.Get(r, "params").(httprouter.Params)
  body := context.Get(r, "body").(*models.HabitResource)

  // Build query w/ account_id to ensure ownership
  q := bson.M{
    "_id": bson.ObjectIdHex(params.ByName("id")),
    "account_id": c.GetCurrentUserId(),
  }

  // Only `name` may be updated
  update := bson.M{
    "$set": bson.M{
      "name": body.Payload.Name,
      "updated_at": lib.MakeTimestamp(),
    },
  }

  repo := models.HabitRepo{c.MongoSession.C("habits")}

  err := repo.Update(q, update)

  if err != nil {
    panic(err)
  }

  w.WriteHeader(204)
  w.Write([]byte("\n"))
}

func DeleteHabit(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("DeleteHabit")

  params := context.Get(r, "params").(httprouter.Params)
  repo := models.HabitRepo{c.MongoSession.C("habits")}

  // Build query w/ account_id to ensure ownership
  q := bson.M{
    "_id": bson.ObjectIdHex(params.ByName("id")),
    "account_id": c.GetCurrentUserId(),
  }

  err := repo.Delete(q)

  if err != nil {
    panic(err)
  }

  w.WriteHeader(204)
  w.Write([]byte("\n"))
}
