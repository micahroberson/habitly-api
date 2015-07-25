package handlers

import (
  // "encoding/json"
  "fmt"
  "net/http"
  "gopkg.in/mgo.v2/bson"
  "github.com/gorilla/context"
  // "github.com/jinzhu/copier"
  "github.com/julienschmidt/httprouter"
  "github.com/micahroberson/habitly-api/lib"
  "github.com/micahroberson/habitly-api/models"
  "time"
)

func boh(t time.Time) time.Time {
    year, month, day := t.Date()
    hour := t.Hour()
    return time.Date(year, month, day, hour, 0, 0, 0, t.Location())
}

func CreateDatapoint(c *lib.AppContext, w http.ResponseWriter, r *http.Request) {
  fmt.Println("CreateDatapoint")

  params := context.Get(r, "params").(httprouter.Params)
  body := context.Get(r, "body").(*models.SecDatapointResource)
  repo := models.DatapointRepo{c.MongoSession.C("datapoints")}

  habitId := bson.ObjectIdHex(params.ByName("habit_id"))

  // Check for existing hourly datapoint
  timestamp := time.Now()
  timestampHour := boh(timestamp)
  min := timestamp.Truncate(time.Minute).Minute()
  sec := timestamp.Truncate(time.Second).Second()

  datapoint, lookupErr := repo.Find(bson.M{
    "habit_id":       habitId,
    "timestamp_hour": timestampHour,
  })

  if lookupErr != nil {
    fmt.Println("Hourly datapoint doesn't exist. Creating a new one...")
    // Create a new datapoint
    // TODO: Move to background worker
    // secDatapoints := [60]models.SecDatapoint{}
    // for i := 0; i < 60; i++ {
    //   secDatapoints[i].Value = 0.0
    // }
    minDatapoints := [60]models.MinDatapoint{}
    for i := 0; i < 60; i++ {
      // copier.Copy(&minDatapoints[i].Values, &secDatapoints)
      minDatapoints[i].Values = [60]models.SecDatapoint{}
    }

    // Set value
    minDatapoints[min].Values[sec].Value = body.Payload.Value

    datapoint = models.Datapoint{
      HabitId:   habitId,
      CreatedAt: lib.MakeTimestamp(),
      UpdatedAt: lib.MakeTimestamp(),
      TimestampHour: timestampHour,
      Values: minDatapoints,
    }
    err := repo.Create(&datapoint)

    if err != nil {
      panic(err)
    }
  } else {
    fmt.Println("Datapoint does exist. Updating now...")
    q := bson.M{
      "_id": datapoint.Id,
    }
    fmt.Println("Datapoint Id: ", datapoint.Id)
    innerSet := bson.M{}
    datapointKey := fmt.Sprintf("values.%d.values.%d.value", min, sec)
    fmt.Println("Datapoint Key: ", datapointKey)
    innerSet[datapointKey] = body.Payload.Value
    update := bson.M{
      "$set": innerSet,
    }

    if err := repo.Update(q, update); err != nil {
      panic(err)
    }
  }

  w.WriteHeader(204)
  w.Write([]byte("\n"))
}