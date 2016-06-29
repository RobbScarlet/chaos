package dao

import (
    "github.com/google/cayley"
    "github.com/google/cayley/quad"
    "errors"
    "strings"
    "chaos/config"
    _ "github.com/google/cayley/graph/bolt"
    "github.com/google/cayley/graph/path"
    "github.com/pquerna/ffjson/ffjson"
    "github.com/google/cayley/graph"
)

//
// AddService 添加服务到Cayley数据库中
// service:服务信息;replaceOrIngore:信息存在覆盖(replace)或者忽略(ingore),非(ingore)值,则都为replace
func AddService(service Service, replaceOrIngore string)(err error){
    // 检测Service数据合法性
    err = checkService(service)
    if err != nil{
        return
    }
    quads, err := convertServiceToQuads(service)
    if err != nil{
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    oldServiceJson := getServiceJson(handle.QuadStore, service.Name)
    if oldServiceJson == ""{ // 不存在该ServiceName对应的服务
        // 添加新的Quad
        return handle.QuadWriter.AddQuadSet(quads)
    }else if strings.ToLower(replaceOrIngore) != "ingore"{ // 替换已有信息
        // TODO 删除旧的服务,并添加新的服务,需要考虑事务
        // 删除旧的Quad
        oldQuads := convertServiceJsonToQuads(service.Name, oldServiceJson)
        for _, oldQuad := range oldQuads{
            err = handle.QuadWriter.RemoveQuad(oldQuad)
            if err != nil{
                return
            }
        }
        // 添加新的Quad
        return handle.QuadWriter.AddQuadSet(quads)
    }

    return
}

//
// GetService 根据服务名获取服务的基本信息
//
func GetService(serviceName string)(service *Service, err error){
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    serviceJson := getServiceJson(handle.QuadStore, serviceName)
    if serviceJson == ""{ // 服务不存在
        err = ServiceNotExist
    }else{
        service = &Service{}
        err = ffjson.Unmarshal([]byte(serviceJson), service)
    }

    return
}

//
// DeleteService 删除服务信息
//
func DeleteService(serviceName string)(service *Service, err error){
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    serviceJson := getServiceJson(handle.QuadStore, serviceName)
    if serviceJson == ""{ // 服务不存在
        err = ServiceNotExist
    }else{
        service = &Service{}
        err = ffjson.Unmarshal([]byte(serviceJson), service)
        if err != nil{
            return
        }
        quads := convertServiceJsonToQuads(serviceName, serviceJson)
        // TODO 此处需要使用事务,必须保证同时删除成功
        for _, quad := range quads{
            err = handle.QuadWriter.RemoveQuad(quad)
            if err != nil{
                return
            }
        }
    }

    return
}

// convertServiceToQuads:将服务数据转化成,Quad切片
func convertServiceToQuads(service Service)(quads []quad.Quad, err error){
    quads = make([]quad.Quad, 0)
    // 服务名称
    quads = append(quads, cayley.Quad(service.Name, ServiceNamePre, service.Name, ""))
    // 服务信息JSON字符串
    serviceJSON, err := ffjson.Marshal(service)
    if err != nil{
        return
    }
    quads = append(quads, cayley.Quad(service.Name, ServiceInfoPre, string(serviceJSON), ""))

    return
}

func convertServiceJsonToQuads(serviceName string, serviceJson string)(quads []quad.Quad){
    quads = make([]quad.Quad, 0)
    // 服务名称
    quads = append(quads, cayley.Quad(serviceName, ServiceNamePre, serviceName, ""))
    // 服务信息JSON字符串
    quads = append(quads, cayley.Quad(serviceName, ServiceInfoPre, serviceJson, ""))

    return
}


// checkService:检测Service数据是否合法
// Service.Name, Service.Category不能为空
func checkService(service Service) (err error){
    if service.Name == ""{
        err = errors.New("服务名称不能为空")
    }
    if service.Category == ""{
        err = errors.New("服务种类不能为空")
    }

    return
}

// getServiceJson:根据服务名获取服务的信息JSON
// serviceName:服务的名称
func getServiceJson(qs graph.QuadStore, serviceName string)(serviceJson string){
    it := path.StartPath(qs, serviceName).Out(ServiceInfoPre).BuildIterator()
    for graph.Next(it){
        serviceJson = qs.NameOf(it.Result())
    }

    return
}
