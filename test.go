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
    //fmt.Println(dao.GetServiceLinkSource("D"))
    //fmt.Println(dao.GetServiceLinkTarget("A"))
    //fmt.Println(dao.GetServiceLinkBoth("E"))

    fmt.Println(dao.GetServiceTree("A", "target"))
    //fmt.Println(dao.GetServiceTreeTarget("B"))
    //fmt.Println(dao.GetServiceTreeTarget("E"))
    //fmt.Println(dao.GetServiceTreeTarget("F"))
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
    AB := &dao.ServiceLink{
        Source: "A",
        Target: "B",
    }
    AE := &dao.ServiceLink{
        Source: "A",
        Target: "E",
    }
    AF := &dao.ServiceLink{
        Source: "A",
        Target: "F",
    }
    BC := &dao.ServiceLink{
        Source: "B",
        Target: "C",
    }
    BD := &dao.ServiceLink{
        Source: "B",
        Target: "D",
    }
    ED := &dao.ServiceLink{
        Source: "E",
        Target: "D",
    }
    FB := &dao.ServiceLink{
        Source: "F",
        Target: "B",
    }

    fmt.Println(dao.AddService(*A, "ignore"))
    fmt.Println(dao.AddService(*B, "ignore"))
    fmt.Println(dao.AddService(*C, "ignore"))
    fmt.Println(dao.AddService(*D, "ignore"))
    fmt.Println(dao.AddService(*E, "ignore"))
    fmt.Println(dao.AddService(*F, "ignore"))
    fmt.Println(dao.AddServiceLink(*AB))
    fmt.Println(dao.AddServiceLink(*AE))
    fmt.Println(dao.AddServiceLink(*AF))
    fmt.Println(dao.AddServiceLink(*BC))
    fmt.Println(dao.AddServiceLink(*BD))
    fmt.Println(dao.AddServiceLink(*ED))
    fmt.Println(dao.AddServiceLink(*FB))

}