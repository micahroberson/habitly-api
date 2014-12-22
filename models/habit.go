package models

import (
  "github.com/codegangsta/martini-contrib/binding"
  "labix.org/v2/mgo/bson"
  "net/http"
  "time"
)

type Habit struct {
  Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name        string        `json:"name" bson:"name"`
  CreatedBy   bson.ObjectId `json:"created_by" bson:"created_by"`
  CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

// This method implements binding.Validator and is executed by the binding.Validate middleware
// Should only be called when creating a new event via a POST request
func (event Habit) Validate(errors *binding.Errors, req *http.Request) {

  if event.Name == "" {
    errors.Fields["name"] = "This field is required"
  } else if len(event.Name) < 5 {
    errors.Fields["name"] = "Too short, minimum 5 characters"
  } else if len(event.Name) > 50 {
    errors.Fields["name"] = "Too long, maximum 50 characters"
  }

  if len(errors.Fields) > 0 {
    errors.Overall["ValidationError"] = "Form validation failed"
  }

}