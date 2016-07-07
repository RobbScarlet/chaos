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
// GetServiceList 根据服务名列表获取服务的基本信息列表
//
func GetServiceList(serviceNameList []string)(serviceList []Service, err error){
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    serviceList = make([]Service, 0)
    if len(serviceNameList) == 0{
        err = ServiceNameListIsNil
        return
    }
    it := path.StartPath(handle.QuadStore, serviceNameList...).Out(ServiceInfoPre).BuildIterator()
    for graph.Next(it){ // 循环获取所有的服务信息JSON
        serviceJson := handle.QuadStore.NameOf(it.Result())
        service := &Service{}
        err = ffjson.Unmarshal([]byte(serviceJson), service)
        if err != nil{
            serviceList = make([]Service, 0)
            return
        }
        serviceList = append(serviceList, *service)
    }
    return
}

//
//  GetAllService 查询所有的服务
//
func GetAllService()(serviceList []Service, err error){
    serviceList = make([]Service, 0)
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    it := path.StartPath(handle.QuadStore).Out(ServiceInfoPre).BuildIterator()
    for graph.Next(it){
        serviceJson := handle.QuadStore.NameOf(it.Result())
        service := &Service{}
        err = ffjson.Unmarshal([]byte(serviceJson), service)
        if err != nil{
            serviceList = make([]Service, 0)
            return
        }
        serviceList = append(serviceList, *service)
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
//  GetServiceLink 获取服务的调用关系
//  linkType:source=调用该服务的所有服务,target=该服务调用的所有服务,both=source+target
func GetServiceLink(serviceName, linkType string)(serviceLinkList []ServiceLink, err error){
    serviceLinkList = make([]ServiceLink, 0)
    if serviceName == ""{
        err = ServiceNameIsNil
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()

    if linkType == "source"{
        serviceLinkList = getServiceLinkSource(serviceName, handle.QuadStore, serviceLinkList)
    }else if linkType == "target"{
        serviceLinkList = getServiceLinkTarget(serviceName, handle.QuadStore, serviceLinkList)
    }else if linkType == "both"{
        // source
        serviceLinkList = getServiceLinkSource(serviceName, handle.QuadStore, serviceLinkList)
        // target
        serviceLinkList = getServiceLinkTarget(serviceName, handle.QuadStore, serviceLinkList)
    }
    return
}

//
//  GetServiceAndServiceLink 获取服务的信息和服务的调用关系
//  linkType:source=调用该服务的所有服务,target=该服务调用的所有服务,both=source+target
func GetServiceAndServiceLink(serviceName, linkType string)(serviceAndServiceLink *ServiceAndServiceLink, err error){
    serviceNameList := make([]string, 0)
    serviceNameMap := make(map[string]string)
    serviceLinkList := make([]ServiceLink, 0)

    if serviceName == ""{
        err = ServiceNameIsNil
        return
    }

    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()
    // 服务调用关系
    if linkType == "source"{
        serviceLinkList = getServiceLinkSource(serviceName, handle.QuadStore, serviceLinkList)
    }else if linkType == "target"{
        serviceLinkList = getServiceLinkTarget(serviceName, handle.QuadStore, serviceLinkList)
    }else if linkType == "both"{
        // source
        serviceLinkList = getServiceLinkSource(serviceName, handle.QuadStore, serviceLinkList)
        // target
        serviceLinkList = getServiceLinkTarget(serviceName, handle.QuadStore, serviceLinkList)
    }
    // 服务的基本信息
    serviceNameList = append(serviceNameList, serviceName)
    serviceNameMap[serviceName] = serviceName
    for _, serviceLink := range serviceLinkList{
        if serviceNameMap[serviceLink.Source] == ""{
            serviceNameList = append(serviceNameList, serviceLink.Source)
            serviceNameMap[serviceLink.Source] = serviceLink.Source
        }
        if serviceNameMap[serviceLink.Target] == ""{
            serviceNameList = append(serviceNameList, serviceLink.Target)
            serviceNameMap[serviceLink.Target] = serviceLink.Target
        }
    }
    serviceList, err := getServiceListByNameList(serviceNameList, handle.QuadStore)
    if err != nil{
        return
    }
    // 服务的基本信息列表和服务调用关系列表
    serviceAndServiceLink = &ServiceAndServiceLink{serviceList, serviceLinkList}

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
//  GetServiceTree 获取该服务的调用关系树,linkType表示调用以serviceName为根,寻找调用该服务和被该服务调用的树
//  linkType:source=调用该服务;target=被该服务调用
func GetServiceTree(serviceName, linkType string)(serviceTree *ServiceTree, err error){
    handle, err := cayley.NewGraph("bolt", config.Bolt_DB_Path, nil)
    if err != nil{
        return
    }
    defer handle.Close()
    // 根节点调用信息
    rootNode := &ServiceTreeNode{ServiceName: serviceName}
    // 树形调用结构中的所有服务
    serviceNameMap := make(map[string]string)
    // 子调用
    rootNode.LinkNodeList = getServiceLinkNode(serviceName, handle.QuadStore, serviceNameMap, linkType)
    // 所有的服务基本信息
    serviceNameList := make([]string, 0)
    for _, serviceName := range serviceNameMap{
        serviceNameList = append(serviceNameList, serviceName)
    }
    serviceList, err := getServiceListByNameList(serviceNameList, handle.QuadStore)
    if err != nil{
        return
    }

    serviceTree = &ServiceTree{
        ServiceList: serviceList,
        RootNode: *rootNode,
    }
    return
}

//  getServiceListByNameList 根据服务名称获取所对应的服务信息
//
func getServiceListByNameList(serviceNameList []string, qs graph.QuadStore)(serviceList []Service, err error){
    serviceList = make([]Service, 0)
    if len(serviceNameList) == 0{
        err = ServiceNameListIsNil
        return
    }
    it := path.StartPath(qs, serviceNameList...).Out(ServiceInfoPre).BuildIterator()
    for graph.Next(it){ // 循环获取所有的服务信息JSON
        serviceJson := qs.NameOf(it.Result())
        service := &Service{}
        err = ffjson.Unmarshal([]byte(serviceJson), service)
        if err != nil{
            serviceList = make([]Service, 0)
            return
        }
        serviceList = append(serviceList, *service)
    }
    return
}

//  GetServiceLinkTarget 获取所有被当前服务调用的g服务
//
func getServiceLinkTarget(serviceName string, qs graph.QuadStore, serviceLinkList []ServiceLink)([]ServiceLink){
    it := path.StartPath(qs, serviceName).Out(ServiceLinkPre).BuildIterator()
    for graph.Next(it){
        target := qs.NameOf(it.Result())
        serviceLink := &ServiceLink{
            Source: serviceName,
            Target: target,
        }
        serviceLinkList = append(serviceLinkList, *serviceLink)
    }

    return serviceLinkList
}

//  GetServiceLinkSource 获取所有调用当前服务的服务
//
func getServiceLinkSource(serviceName string, qs graph.QuadStore, serviceLinkList []ServiceLink)([]ServiceLink){
    it := path.StartPath(qs, serviceName).In(ServiceLinkPre).BuildIterator()
    for graph.Next(it){
        source := qs.NameOf(it.Result())
        serviceLink := &ServiceLink{
            Source: source,
            Target: serviceName,
        }
        serviceLinkList = append(serviceLinkList, *serviceLink)
    }

    return serviceLinkList
}

// getServiceLinkTargetNode 获取该服务调用的服务/被调用该服务的服务列表
// linkType:source=调用该服务;target=被该服务调用
func getServiceLinkNode(serviceName string, qs graph.QuadStore, serviceNameMap map[string]string, linkType string)(linkNodeList []ServiceTreeNode){
    serviceNameMap[serviceName] = serviceName
    linkNodeList = make([]ServiceTreeNode, 0)
    var it graph.Iterator
    if linkType == "source"{
        it = path.StartPath(qs, serviceName).In(ServiceLinkPre).BuildIterator()
    }else if linkType == "target"{
        it = path.StartPath(qs, serviceName).Out(ServiceLinkPre).BuildIterator()
    }else{
        return
    }

    for graph.Next(it){
        targetServiceName := qs.NameOf(it.Result())
        linkNode := &ServiceTreeNode{ServiceName: targetServiceName}
        // 如果没有计算过子调用
        if serviceNameMap[targetServiceName] == ""{
            linkNode.LinkNodeList = getServiceLinkNode(targetServiceName, qs, serviceNameMap, linkType)
        }

        linkNodeList = append(linkNodeList, *linkNode)
    }
    return
}

// convertServiceToQuads:将服务数据转化成Quad列表
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
// convertServiceJsonToQuads:将服务数据转化成Quad列表
func convertServiceJsonToQuads(serviceName string, serviceJson string)(quads []quad.Quad){
    quads = make([]quad.Quad, 0)
    // 服务名称
    quads = append(quads, cayley.Quad(serviceName, ServiceNamePre, serviceName, ""))
    // 服务信息JSON字符串
    quads = append(quads, cayley.Quad(serviceName, ServiceInfoPre, serviceJson, ""))

    return
}

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