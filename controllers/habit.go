package controllers

import (
  "github.com/codegangsta/martini-contrib/render"
  "github.com/go-martini/martini"
  "github.com/micahroberson/habitly-api/models"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
)

// Called on a GET to /habits
// Returns a list of all habits
func GetAllHabits(db *mgo.Database, r render.Render) {

  var habits []models.Habit

  err := db.C("habits").Find(nil).Sort("created_at").All(&habits)

  if err != nil {
    panic(err)
  }

  r.JSON(200, habits)
}

// Called on a Get to /habits/:id
// Returns a single habit with the id :id
func GetHabit(db *mgo.Database, r render.Render, p martini.Params) {

  var id bson.ObjectId
  var habit models.Habit

  if bson.IsObjectIdHex(p["id"]) {
    id = bson.ObjectIdHex(p["id"])
  } else {
    r.JSON(400, "Bad Request: Invalid ObjectId")
    return
  }

  err := db.C("habits").FindId(id).One(&habit)

  // TODO Check for 404

  if err != nil {
    panic(err)
  }

  r.JSON(200, habit)

}

// Called on a POST to /habits
// Assuming valid habit; adds the given habit
func AddHabit(db *mgo.Database, r render.Render, habit models.Habit) {

  // Create a unique id
  habit.Id = bson.NewObjectId()

  // TODO Should be the user Id
  habit.CreatedBy = bson.NewObjectId()
  habit.CreatedAt = time.Now().UTC()

  err := db.C("habits").Insert(habit)
  if err != nil {
    panic(err)
  }

  r.JSON(201, habit)
}