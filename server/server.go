package server

import (
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/codegangsta/martini-contrib/cors"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/go-martini/martini"
  // "github.com/martini-contrib/auth"
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

  var LoginRequired martini.Handler = func() martini.Handler {
    return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
      
    }
  }

  m.Group("/api/v1", func(r martini.Router) {
    r.Get("/habits", controllers.GetAllHabits)
    r.Get("/habits/:id", controllers.GetHabit)
    r.Post("/habits", binding.Json(models.Habit{}), binding.ErrorHandler, controllers.AddHabit)
  })

  m.Group("/api/v1", func(r martini.Router) {
    r.Post("/users", binding.Json(models.User{}), binding.ErrorHandler, controllers.RegisterUser)
    r.Post("/login", binding.Json(models.User{}), binding.ErrorHandler, controllers.LoginUser)
    // r.Delete("/logout", controllers.UserLogout)
  })

  // Setup comment routes
  // m.Get(`/habits/:habit_id/comments`, controllers.GetAllComments)
  // m.Post(`/habits/:habit_id/comments`, binding.Json(models.Comment{}), binding.ErrorHandler, controllers.AddComment)

  // TODO Update, Delete for habits
  //m.Put(`/habits/:id`, UpdateHabit)
  //m.Delete(`/habits/:id`, DeleteHabit)

  // Add the router action

  return m
}