package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
    "net/http"
)

type serviceHandler struct{

}
func NewServiceHandler() *serviceHandler{
    return &serviceHandler{}
}

func (handler *serviceHandler)findAllServices(request *restful.Request, response *restful.Response){

}

func (handler *serviceHandler)findService(request *restful.Request, response *restful.Response){
    serviceName := request.PathParameter("service-name")
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
    }
    service, err := dao.GetService(serviceName)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    response.WriteEntity(service)
}

func (handler *serviceHandler)createService(request *restful.Request, response *restful.Response){
    service := dao.Service{Name: request.PathParameter("service-name")}
    err := request.ReadEntity(service)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    err = dao.AddService(service, "ingore")
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    response.WriteHeader(http.StatusCreated)
}

func (handler *serviceHandler)removeService(request *restful.Request, response *restful.Response){

}

func (handler *serviceHandler)createServiceLink(request *restful.Request, response *restful.Response){
    serviceLink := dao.ServiceLink{}
    err := request.ReadEntity(serviceLink)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    err = dao.AddServiceLink(serviceLink)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    response.WriteHeader(http.StatusCreated)
}

func (handler *serviceHandler)findServiceLink(request *restful.Request, response *restful.Response){
    serviceName := request.PathParameter("service-name")
    linkType := request.QueryParameter("linkType")
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
    }

    // 已经查询过服务信息的服务
    serviceMap := make(map[string]string)
    // 服务信息
    serviceList := make([]dao.Service, 0)
    // 服务调用关系列表
    serviceLinkList, err := getServiceLinkByType(serviceName, linkType)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }

    for _, serviceLink := range serviceLinkList{
        // Source
        serviceList = getServiceByName(serviceLink.Source, serviceMap, serviceList)
        // Target
        serviceList = getServiceByName(serviceLink.Target, serviceMap, serviceList)
    }
    serviceList = getServiceByName(serviceName, serviceMap, serviceList)

    response.WriteEntity(&ServiceAndServiceLink{serviceList, serviceLinkList})
}

func getServiceLinkByType(serviceName, linkType string)(serviceLinkList []dao.ServiceLink, err error){
    if linkType == "source"{
        serviceLinkList, err = dao.GetServiceLinkSource(serviceName)
    }else if linkType == "target"{
        serviceLinkList, err = dao.GetServiceLinkTarget(serviceName)
    }else{
        serviceLinkList, err = dao.GetServiceLinkBoth(serviceName)
    }
    return
}

func getServiceByName(servieName string, existServiceMap map[string]string, serviceList []dao.Service)([]dao.Service){
    if existServiceMap[servieName] == ""{
        service, err := dao.GetService(servieName)
        if err == nil{
            serviceList = append(serviceList, *service)
            existServiceMap[servieName] = servieName
        }
    }

    return serviceList
}
