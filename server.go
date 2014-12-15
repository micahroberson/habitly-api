package server

import (
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/codegangsta/martini-contrib/cors"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/go-martini/martini"
  goauth2 "github.com/golang/oauth2"
  "github.com/martini-contrib/oauth2"
  "github.com/martini-contrib/sessions"
  "github.com/micahroberson/habitly-api/config"
  "github.com/micahroberson/habitly-api/controllers"
  "github.com/micahroberson/habitly-api/db"
  "github.com/micahroberson/habitly-api/models"
)

func NewServer(databaseName string) *martini.ClassicMartini {

  m := martini.Classic()
  c := config.GetConfig()

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

  // Google OAuth
  m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte(c.Cookie_Auth),
    []byte(c.Cookie_Enc))))

  m.Use(oauth2.Google(
    goauth2.Client(c.Client_Id, c.Client_Secret),
    goauth2.RedirectURL(c.OAuth_Callback),
    goauth2.Scope("email"),
  ))

  // Setup event routes
  m.Get(`/habits`, controllers.GetAllHabits)
  m.Get(`/habits/:id`, controllers.GetHabit)
  m.Post(`/habits`, binding.Json(models.Habit{}), binding.ErrorHandler, controllers.AddHabit)

  // Setup comment routes
  // m.Get(`/habits/:habit_id/comments`, controllers.GetAllComments)
  // m.Post(`/habits/:habit_id/comments`, binding.Json(models.Comment{}), binding.ErrorHandler, controllers.AddComment)

  // User route for Oauth
  m.Get(`/users`,  oauth2.LoginRequired, controllers.GetLoggedInUser)

  // TODO Update, Delete for habits
  //m.Put(`/habits/:id`, UpdateHabit)
  //m.Delete(`/habits/:id`, DeleteHabit)

  // Add the router action

  return m
}