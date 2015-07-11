package main

import (
  "crypto/rsa"
  "encoding/json"
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "os"
  "reflect"
  "time"
  "github.com/gorilla/context"
  // "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/justinas/alice"
  "github.com/micahroberson/habitly-api/handlers"
  "github.com/micahroberson/habitly-api/lib"
  "github.com/micahroberson/habitly-api/models"
  "github.com/micahroberson/habitly-api/router"
)

const (
  privateKeyPath = "keys/habitly-api.rsa"
  publicKeyPath  = "keys/habitly-api.rsa.pub"
)

var (
  verifyKey *rsa.PublicKey
  signKey   *rsa.PrivateKey
)

func envOrDefault(key_name, default_value string) (env_value string) {
  if e := os.Getenv(key_name); len(e) == 0 {
    env_value = default_value
  } else {
    env_value = e
  }
  return
}

var (
  listenPort    = "8080"
  mongouser     = envOrDefault("MONGOUSER", "")
  mongosecret   = envOrDefault("MONGOSECRET", "")
  mongohost     = envOrDefault("MONGOHOST", "localhost")
  mongodb       = envOrDefault("MONGODB", "habitly-api")
)

type Errors struct {
  Errors []*Error `json:"errors"`
}

type Error struct {
  Id     string `json:"id"`
  Status int    `json:"status"`
  Title  string `json:"title"`
  Detail string `json:"detail"`
}

func WriteError(w http.ResponseWriter, err *Error) {
  w.Header().Set("Content-Type", "application/vnd.api+json")
  w.WriteHeader(err.Status)
  json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

var (
  ErrBadRequest           = &Error{"bad_request", 400, "Bad request", "Request body is not well-formed. It must be JSON."}
  ErrNotAuthorized        = &Error{"not_authorized", 401, "Not Authorized", "You are not authorized to access this resource."}
  ErrNotAcceptable        = &Error{"not_acceptable", 406, "Not Acceptable", "Accept header must be set to 'application/vnd.api+json'."}
  ErrUnsupportedMediaType = &Error{"unsupported_media_type", 415, "Unsupported Media Type", "Content-Type header must be set to: 'application/vnd.api+json'."}
  ErrInternalServer       = &Error{"internal_server_error", 500, "Internal Server Error", "Something went wrong."}
)

func appHandler(c *lib.AppContext, h func(*lib.AppContext, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
  fn := func(w http.ResponseWriter, r *http.Request) {
    h(c, w, r)
  }
  return fn
}

func authenticationHandler(c *lib.AppContext) func(http.Handler) http.Handler {
  m := func(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
      // Get jwt from header
      tokenString := r.Header.Get("X-Access-Token")
      if tokenString == "" {
        WriteError(w, ErrNotAuthorized)
        return
      }

      token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return verifyKey, nil
      })
      switch err.(type) {
        case nil:
          if !token.Valid {
            WriteError(w, ErrNotAuthorized)
            return
          }
          // Set CurrentUserId
          // fmt.Println("Token: ", token)
          c.SetCurrentUserId(token.Claims["user_id"].(string))
        case *jwt.ValidationError:
          vErr := err.(*jwt.ValidationError)
          switch vErr.Errors {
          case jwt.ValidationErrorExpired:
            WriteError(w, ErrNotAuthorized)
            return
          default:
            WriteError(w, ErrInternalServer)
            return
          }
        default:
          WriteError(w, ErrNotAuthorized)
          return
      }

      next.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn)
  }

  return m
}

func recoverHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    defer func() {
      if err := recover(); err != nil {
        log.Printf("panic: %+v", err)
        WriteError(w, ErrInternalServer)
      }
    }()

    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    t1 := time.Now()
    next.ServeHTTP(w, r)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
  }

  return http.HandlerFunc(fn)
}

func acceptHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Accept") != "application/vnd.api+json" {
      WriteError(w, ErrNotAcceptable)
      return
    }

    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func contentTypeHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
      if r.Header.Get("Content-Type") != "application/vnd.api+json" {
        WriteError(w, ErrUnsupportedMediaType)
        return
      }

      next.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn)
}

func bodyHandler(v interface{}) func(http.Handler) http.Handler {
  t := reflect.TypeOf(v)

  m := func(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
      val := reflect.New(t).Interface()
      err := json.NewDecoder(r.Body).Decode(val)

      if err != nil {
        WriteError(w, ErrBadRequest)
        return
      }

      if next != nil {
        context.Set(r, "body", val)
        next.ServeHTTP(w, r)
      }
    }

    return http.HandlerFunc(fn)
  }

  return m
}

func init() {
  signBytes, err := ioutil.ReadFile(privateKeyPath)
    fatal(err)

    signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
    fatal(err)

    verifyBytes, err := ioutil.ReadFile(publicKeyPath)
    fatal(err)

    verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
    fatal(err)
}

func fatal(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  fmt.Println("Starting up habitly-api...")

  mongoAuthString := "mongodb://"
  if len(mongouser) != 0 || len(mongosecret) != 0 {
    mongoAuthString = mongoAuthString + mongouser + ":" + mongosecret + "@"
  }

  mongoConnectionString := mongoAuthString + mongohost + "/" + mongodb
  mgoSession, err := mgo.Dial(mongoConnectionString)

  if err != nil {
    fmt.Println("Could not connect to MongoDB", mongoConnectionString)
    panic(err)
  }
  defer mgoSession.Close()
  mgoSession.SetMode(mgo.Monotonic, true)

  appContext := &lib.AppContext{
    MongoSession: mgoSession.DB(""),
    VerifyKey: verifyKey,
    SignKey: signKey,
  }

  commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler)
  router := router.NewRouter()
  router.Post("/api/users/sign-in", commonHandlers.Append(contentTypeHandler, bodyHandler(models.UserResource{})).ThenFunc(appHandler(appContext, handlers.AuthenticateUser)))
  router.Post("/api/users/register", commonHandlers.Append(contentTypeHandler, bodyHandler(models.UserResource{})).ThenFunc(appHandler(appContext, handlers.RegisterUser)))
  router.Get("/api/habits", commonHandlers.Append(authenticationHandler(appContext)).ThenFunc(appHandler(appContext, handlers.GetAllHabits)))
  router.Get("/api/habits/:id", commonHandlers.Append(authenticationHandler(appContext)).ThenFunc(appHandler(appContext, handlers.GetHabit)))
  router.Post("/api/habits", commonHandlers.Append(contentTypeHandler, authenticationHandler(appContext), bodyHandler(models.HabitResource{})).ThenFunc(appHandler(appContext, handlers.CreateHabit)))
  router.Put("/api/habits/:id", commonHandlers.Append(contentTypeHandler, authenticationHandler(appContext), bodyHandler(models.HabitResource{})).ThenFunc(appHandler(appContext, handlers.UpdateHabit)))
  router.Delete("/api/habits/:id", commonHandlers.Append(authenticationHandler(appContext)).ThenFunc(appHandler(appContext, handlers.DeleteHabit)))


  // router.
  //   Methods("POST").
  //   Path("/api/users/sign-in").
  //   Name("Authenticate").
  //   Handler(appHandler{false, context, handlers.Authenticate})

  // router.
  //   Methods("POST").
  //   Path("/api/users/register").
  //   Name("RegisterUser").
  //   Handler(appHandler{false, context, handlers.RegisterUser})

  // router.
  //   Methods("GET").
  //   Path("/api/habits").
  //   Name("HabitsIndex").
  //   Handler(appHandler{true, context, handlers.GetAllHabits})

  // http.Handle("/", router)

  fmt.Printf("habitly-api server listening on port: %s...\n\n", listenPort)
  http.ListenAndServe(fmt.Sprintf(":%s", listenPort), router)
}