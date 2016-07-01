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
// service:服务信息;replaceOrIngore:信息存在覆盖(replace)或者忽略(ignore),非(ignore)值,则都为replace
func AddService(service Service, replaceOrIgnore string)(err error){
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
    }else if strings.ToLower(replaceOrIgnore) != "ignore"{ // 替换已有信息
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
// GetService 查询所有的服务
//
//func GetAllService()(service []Service, err error){
//    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
//    if err != nil{
//        return
//    }
//    defer handle.Close()
//
//
//}

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
//
//  AddServiceLink 添加服务调用关系
//
func AddServiceLink(serviceLink ServiceLink)(err error){
    quad, err := convertServiceLinkToQuad(serviceLink)
    if err != nil{
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    err = handle.QuadWriter.AddQuad(quad)
    if err == graph.ErrQuadExists{
        return nil
    }
    return
}

//
//  AddServiceLink 添加服务调用关系列表
//
func AddServiceLinkList(serviceLinkList ServiceLinkList)(err error){
    quads, err := convertServiceLinkListToQuads(serviceLinkList)
    if err != nil{
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()
    // TODO 增加事务
    for _, quad := range quads{
        err = handle.QuadWriter.AddQuad(quad)
        if err == graph.ErrQuadExists{
            continue
        }else{
            return
        }
    }

    return
}

//
//  DeleteServiceLink 删除服务调用关系
//
func DeleteServiceLink(serviceLink ServiceLink)(err error){
    quad, err := convertServiceLinkToQuad(serviceLink)
    if err != nil{
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    return handle.QuadWriter.RemoveQuad(quad)
}

//
//  GetServiceLinkTarget 获取所有被当前服务调用的服务
//
func GetServiceLinkTarget(serviceName string)(serviceLinkList []ServiceLink, err error){
    if serviceName == ""{
        err = errors.New("查询服务调用关系的服务名称不能为空")
        return
    }
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    serviceLinkList = make([]ServiceLink, 0)
    it := path.StartPath(handle.QuadStore, serviceName).Out(ServiceLinkPre).BuildIterator()
    for graph.Next(it){
        target := handle.QuadStore.NameOf(it.Result())
        serviceLink := &ServiceLink{
            Source: serviceName,
            Target: target,
        }
        serviceLinkList = append(serviceLinkList, *serviceLink)
    }
    return
}

//
//  GetServiceLinkSource 获取所有调用当前服务的服务
//
func GetServiceLinkSource(serviceName string)(serviceLinkList []ServiceLink, err error){
    if serviceName == ""{
        err = errors.New("查询服务调用关系的服务名称不能为空")
        return
    }
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    serviceLinkList = make([]ServiceLink, 0)
    it := path.StartPath(handle.QuadStore, serviceName).In(ServiceLinkPre).BuildIterator()
    for graph.Next(it){
        source := handle.QuadStore.NameOf(it.Result())
        serviceLink := &ServiceLink{
            Source: source,
            Target: serviceName,
        }
        serviceLinkList = append(serviceLinkList, *serviceLink)
    }
    return
}

//
//  GetServiceLinkBoth 获取所有被当前服务和调用当前服务的服务
//
func GetServiceLinkBoth(serviceName string)(serviceLinkList []ServiceLink, err error){
    serviceLinkList = make([]ServiceLink, 0)

    // 被该服务调用的服务
    targetServiceLinks, err := GetServiceLinkTarget(serviceName)
    if err != nil{
        return
    }else{
        for _, targetServiceLink := range targetServiceLinks{
            serviceLinkList = append(serviceLinkList, targetServiceLink)
        }
    }
    // 调用该服务的服务
    sourceServiceLinks, err := GetServiceLinkSource(serviceName)
    if err != nil{
        return
    }else{
        for _, sourceServiceLink := range sourceServiceLinks{
            serviceLinkList = append(serviceLinkList, sourceServiceLink)
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
    }else if service.Category == ""{
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
// getServiceJson:获取所有服务的信息JSON
//func getAllServiceJson(qs graph.QuadStore)(serviceJson string){
//    it := path.StartPath(qs)
//}

// checkServiceLink 校验服务调用关系
// 调用方和被调用方不能为空
func convertServiceLinkToQuad(serviceLink ServiceLink)(quad quad.Quad, err error){
    if serviceLink.Source == ""{
        err = errors.New("服务调用关系,调用服务名不能为空")
    }else if serviceLink.Target == ""{
        err = errors.New("服务调用关系,被调用服务名不能为空")
    }else{
        quad = cayley.Quad(serviceLink.Source, ServiceLinkPre, serviceLink.Target, "")
    }

    return
}

// convertServiceLinkListToQuads
// 调用方和被调用方列表不能为空
func convertServiceLinkListToQuads(serviceLinkList ServiceLinkList)(quads []quad.Quad, err error){
    if serviceLinkList.Source == ""{
        err = errors.New("服务调用关系,调用服务名不能为空")
    }else if(len(serviceLinkList.TargetList) == 0){
        err = errors.New("服务调用关系,被调用服务名列表不能为空")
    }else{
        quads = make([]quad.Quad, 0)

        targetList := serviceLinkList.TargetList
        for _, target := range targetList{
            if target == ""{
                err = errors.New("服务调用关系,被调用服务名不能为空")
                break
            }else{
                quads = append(quads, cayley.Quad(serviceLinkList.Source, ServiceLinkPre, target, ""))
            }
        }
    }

    return
}