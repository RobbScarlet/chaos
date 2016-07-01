package main

import (
    _ "github.com/google/cayley/graph/bolt"
    "github.com/google/cayley/graph"
    "chaos/config"
    "chaos/dao"
    "fmt"
)

func main() {
    //initServiceAndLink()
    fmt.Println(dao.GetServiceLinkSource("D"))
    fmt.Println(dao.GetServiceLinkTarget("A"))
    fmt.Println(dao.GetServiceLinkBoth("B"))
}

func initServiceAndLink(){
    err := graph.InitQuadStore("bolt", config.Bolt_DB_Path, nil)
    fmt.Println(err)

    A := &dao.Service{
        Name: "A",
        Category: "Application",
        Description: "Service-A",
    }
    B := &dao.Service{
        Name: "B",
        Category: "Application",
        Description: "Service-B",
    }
    C := &dao.Service{
        Name: "C",
        Category: "Base",
        Description: "Service-C",
    }
    D := &dao.Service{
        Name: "D",
        Category: "Application",
        Description: "Service-D",
    }
    E := &dao.Service{
        Name: "E",
        Category: "Application",
        Description: "Service-E",
    }
    F := &dao.Service{
        Name: "F",
        Category: "Base",
        Description: "Service-F",
    }
    ABEF := &dao.ServiceLinkList{
        Source: "A",
        TargetList: []string{"B", "E", "F"},
    }
    BCD := &dao.ServiceLinkList{
        Source: "B",
        TargetList: []string{"C", "D"},
    }
    ED := &dao.ServiceLink{
        Source: "E",
        Target: "D",
    }

    fmt.Println(dao.AddService(*A, "ignore"))
    fmt.Println(dao.AddService(*B, "ignore"))
    fmt.Println(dao.AddService(*C, "ignore"))
    fmt.Println(dao.AddService(*D, "ignore"))
    fmt.Println(dao.AddService(*E, "ignore"))
    fmt.Println(dao.AddService(*F, "ignore"))
    fmt.Println(dao.AddServiceLinkList(*ABEF))
    fmt.Println(dao.AddServiceLinkList(*BCD))
    fmt.Println(dao.AddServiceLink(*ED))
}