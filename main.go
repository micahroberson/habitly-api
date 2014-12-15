package main

import (
  "github.com/micahroberson/habitly-api/server/server"
)

func main() {
  //Create a new server object and run it
  server := server.NewServer("habitly-api")
  server.Run()
}