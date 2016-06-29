package main

import (
    _ "github.com/google/cayley/graph/bolt"
    "github.com/google/cayley/graph"
    "chaos/config"
    "chaos/dao"
    "fmt"
)

func main() {
    err := graph.InitQuadStore("bolt", config.Bolt_DB_Path, nil)
    fmt.Println(err)

    service1 := dao.Service{
        Name: "zs-service",
        Category: "Application",
        FullName: "zs-thrift-service1",
    }

    service2 := dao.Service{
        Name: "zs-service",
        Category: "Application",
        FullName: "zs-thrift-service2",
    }

    fmt.Println(dao.AddService(service1, "replace"))
    fmt.Println(dao.GetService("zs-service"))
    fmt.Println(dao.AddService(service2, "ingore"))
    fmt.Println(dao.GetService("zs-service"))
    fmt.Println(dao.AddService(service2, "replace"))
    fmt.Println(dao.GetService("zs-service"))
    fmt.Println(dao.DeleteService("zhaoshang-service"))
    fmt.Println(dao.GetService("zs-service"))
}