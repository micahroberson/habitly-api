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
  "time"
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
  habitId := bson.ObjectIdHex(params.ByName("id"))

  q := bson.M{
    "account_id": c.GetCurrentUserId(),
    "_id":        habitId,
  }

  repo := models.HabitRepo{c.MongoSession.C("habits")}
  habit, err := repo.Find(q)

  if err != nil {
    panic(err)
  }

  // TODO: Use current month for past 30 day aggregation
  timestampHourStart, err := time.Parse(time.RFC3339Nano, "2015-07-01T00:00:00Z")
  if err != nil {
    panic(err)
  }

  pipeline := []bson.M{
    bson.M{
      "$match": bson.M{
        "habit_id": habitId,
        "timestamp_hour": bson.M{"$gte": timestampHourStart},
      },
    },
    bson.M{"$unwind": "$values"},
    bson.M{"$unwind": "$values.values"},
    bson.M{
      "$match": bson.M{
        "values.values.value": bson.M{"$gt": 0},
      },
    },
    bson.M{
      "$project": bson.M{
        "value": "$values.values.value",
        "habit_id": 1,
        "timestamp_hour": 1,
      },
    },
    bson.M{
      "$group": bson.M{
        "_id": bson.M{
          "day": bson.M{"$dayOfMonth": "$timestamp_hour"},
          "month": bson.M{"$month": "$timestamp_hour"},
          "year": bson.M{"$year": "$timestamp_hour"},
        },
        "total": bson.M{"$sum": "$value"},
      },
    },
  }

  datapointsRepo := models.DatapointRepo{c.MongoSession.C("datapoints")}
  datapoints, err := datapointsRepo.Aggregate(pipeline)

  if err != nil {
    panic(err)
  }

  habit.Payload.Datapoints = datapoints

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
