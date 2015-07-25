package models

import (
  "fmt"
  "time"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type MinDatapoint struct {
  Values           [60]SecDatapoint     `json:"values"`
}

type SecDatapoint struct {
  Value            float64             `json:"value" bson:"value"`
}

type Datapoint struct {
  Id               bson.ObjectId       `json:"id,omitempty" bson:"_id,omitempty"`
  HabitId          bson.ObjectId       `json:"habit_id" bson:"habit_id"`
  CreatedAt        int64               `json:"created_at" bson:"created_at"`
  UpdatedAt        int64               `json:"updated_at" bson:"updated_at"`
  TimestampHour    time.Time           `json:"timestamp_hour" bson:"timestamp_hour"`
  Values           [60]MinDatapoint    `json:"values"`
}

type SecDatapointResource struct {
  Payload SecDatapoint `json:"payload"`
}

type DatapointCollection struct {
  Payload []Datapoint `json:"payload"`
}

type DatapointRepo struct {
  Coll *mgo.Collection
}

type Day struct {
  Day               int64               `json:"day" bson:"day"`
  Month             int64               `json:"month" bson:"month"`
  Year              int64               `json:"year" bson:"year"`
}

type DatapointAggregation struct {
  Id                Day                 `json:"-" bson:"_id"`
  Date              int64               `json:"date" bson:"date"`
  Total             int64               `json:"total" bson:"total"`
}

type DatapointAggregationCollection struct {
  Payload []DatapointAggregation `json:"payload"`
}

// func (r *DatapointRepo) All(q bson.M) (DatapointCollection, error) {
//   result := DatapointCollection{[]Datapoint{}}
//   err := r.Coll.Find(q).Sort("created_at").All(&result.Payload)
//   if err != nil {
//     return result, err
//   }
//   return result, nil
// }

func (r *DatapointRepo) Aggregate(pipeline []bson.M) ([]DatapointAggregation, error) {
  result := []DatapointAggregation{}
  err := r.Coll.Pipe(pipeline).All(&result)
  if err != nil {
    return result, err
  }
  for i := range result {
    datapoint := result[i]
    dateString := fmt.Sprintf("%02d-%02d-%d", datapoint.Id.Month, datapoint.Id.Day, datapoint.Id.Year)
    date, err := time.Parse("01-02-2006", dateString)
    if err != nil {
      return result, err
    }
    result[i].Date = (date.UnixNano() / int64(time.Millisecond))
  }
  return result, nil
}

func (r *DatapointRepo) Find(q bson.M) (Datapoint, error) {
  result := Datapoint{}
  err := r.Coll.Find(q).One(&result)
  if err != nil {
    return result, err
  }
  return result, nil
}

func (r *DatapointRepo) Create(datapoint *Datapoint) error {
  id := bson.NewObjectId()
  _, err := r.Coll.UpsertId(id, datapoint)
  if err != nil {
    return err
  }
  datapoint.Id = id

  return nil
}

func (r *DatapointRepo) Update(q bson.M, update bson.M) error {
  err := r.Coll.Update(q, update)
  if err != nil {
    return err
  }

  return nil
}
