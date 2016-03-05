package main

import (
	"github.com/ManiacMike/gwork"
  // "fmt"
  "time"
  "strconv"
  "net/http"
)

func StaticServer(w http.ResponseWriter, req *http.Request) {
	// http.ServeFile(w, req, "Web/index.html")
	staticHandler := http.FileServer(http.Dir("Web"))
	staticHandler.ServeHTTP(w, req)
	return
}

func main() {

  http.HandleFunc("/demo", StaticServer)


  gwork.SetGenerateUid(func()string{
    id := int(time.Now().Unix())
    return strconv.Itoa(id)
  });

  gwork.SetGetConnCallback(func(uid string,room *gwork.Room){
    welcome := map[string]interface{}{
      "type" : "welcome",
      "id" : uid,
    }
    room.Broadcast(gwork.JsonEncode(welcome))
  })

  gwork.SetLoseConnCallback(func(uid string,room *gwork.Room){
    close := map[string]interface{}{
      "type" : "close",
      "id" : uid,
    }
    room.Broadcast(gwork.JsonEncode(close))
  })

  gwork.Init(func(receiveNodes map[string]interface{}, uid string, room *gwork.Room){
    receiveType := receiveNodes["type"]
    if receiveType == "login" {

    } else if receiveType == "update" {
      var name interface{}
      var ok bool
      if name,ok = receiveNodes["name"];ok == false{
        name = "Guest."+uid
      }
      reply := map[string]interface{}{
        "type" : "update",
        "id"  : uid,
        "angle" : receiveNodes["angle"].(float64),
        "momentum" : receiveNodes["momentum"].(float64),
        "x" :  receiveNodes["x"].(float64),
        "y" :  receiveNodes["y"].(float64),
        "life" : 1,
        "name" : name,
        "authorized" : false,
      }
      room.Broadcast(gwork.JsonEncode(reply))
    } else if receiveType == "message" {
      reply := map[string]interface{}{
        "type" : "message",
        "id" : uid,
        "message" : receiveNodes["message"].(string),
      }
      room.Broadcast(gwork.JsonEncode(reply))
    }
  })
}
