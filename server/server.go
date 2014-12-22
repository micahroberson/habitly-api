package server

import (
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/codegangsta/martini-contrib/cors"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/go-martini/martini"
  "github.com/micahroberson/habitly-api/controllers"
  "github.com/micahroberson/habitly-api/db"
  "github.com/micahroberson/habitly-api/models"
)

func NewServer(databaseName string) *martini.ClassicMartini {

  m := martini.Classic()

  // Setup middleware
  m.Use(db.DB(databaseName))
  m.Use(render.Renderer())
  m.Use(cors.Allow(&cors.Options{
    AllowOrigins:     []string{"http://localhost*"},
    AllowMethods:     []string{"POST", "GET"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
  }))

  // Setup event routes
  m.Get(`/habits`, controllers.GetAllHabits)
  m.Get(`/habits/:id`, controllers.GetHabit)
  m.Post(`/habits`, binding.Json(models.Habit{}), binding.ErrorHandler, controllers.AddHabit)

  // Setup comment routes
  // m.Get(`/habits/:habit_id/comments`, controllers.GetAllComments)
  // m.Post(`/habits/:habit_id/comments`, binding.Json(models.Comment{}), binding.ErrorHandler, controllers.AddComment)

  // TODO Update, Delete for habits
  //m.Put(`/habits/:id`, UpdateHabit)
  //m.Delete(`/habits/:id`, DeleteHabit)

  // Add the router action

  return m
}