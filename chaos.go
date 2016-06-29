package main

import (
    "chaos/config"
    "github.com/google/cayley/graph"
    "net/http"
    "github.com/emicklei/go-restful"
    "chaos/api"
    "log"
    "fmt"
)

func main() {
    err := graph.InitQuadStore("bolt", config.Bolt_DB_Path, nil)
    if err != nil && err != graph.ErrDatabaseExists {
        panic(err)
    }

    container := restful.NewContainer()
    err = api.Register(container)
    if err != nil{
        panic(err)
    }
    log.Printf("start listening on localhost:10086")
    server := &http.Server{Addr: ":10086", Handler: container}
    fmt.Println(server.ListenAndServe())
}
