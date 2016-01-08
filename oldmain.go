package main

import (
  "log"
  "os"
  "net/http"
)

//TODO User spf13/Cobra for CLI
//TODO Check error handling for GO

func pingRegistry()  {
  response, err := http.Get("http://localhost:5000/v2")

  if err != nil {
    log.Println("Error: ",err)
    os.Exit(1)
  }
  log.Println("Ping Registry: Status %d",response.StatusCode)
}


func main() {
  log.Println("Hello World!")
  pingRegistry()
  os.Exit(0)
}
