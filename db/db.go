package db

import (
  "log"
  "github.com/go-martini/martini"
  "labix.org/v2/mgo"
)

// DB Returns a martini.Handler
func DB(databaseName string) martini.Handler {
  session, err := mgo.Dial("mongodb://localhost")
  if err != nil {
    log.Println("Could not contact mongodb on localhost")
    panic(err)
  }

  return func(c martini.Context) {
    s := session.Clone()
    c.Map(s.DB(databaseName))

    index := mgo.Index{
      Key: []string{"email"},
      Unique: true,
      DropDups: true,
      Background: false,
      Sparse: true,
    }
    s.DB(databaseName).C("users").EnsureIndex(index)

    defer s.Close()
    c.Next()
  }
}