package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
    "net/http"
    "strings"
)

type serviceHandler struct{

}
func NewServiceHandler() *serviceHandler{
    return &serviceHandler{}
}

func (handler *serviceHandler)findAllServices(request *restful.Request, response *restful.Response){
    serviceList, err := dao.GetAllService()
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.WriteEntity(serviceList)
}

func (handler *serviceHandler)findService(request *restful.Request, response *restful.Response){
    serviceName := request.PathParameter("service-name")
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
        return
    }
    service, err := dao.GetService(serviceName)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.WriteEntity(service)
}

func (handler *serviceHandler)createService(request *restful.Request, response *restful.Response){
    service := &dao.Service{Name: request.PathParameter("service-name")}
    err := request.ReadEntity(service)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    err = dao.AddService(*service, "replace")
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.WriteHeader(http.StatusCreated)
}

func (handler *serviceHandler)removeService(request *restful.Request, response *restful.Response){

}

func (handler *serviceHandler)createServiceLink(request *restful.Request, response *restful.Response){
    serviceLink := &dao.ServiceLink{}
    err := request.ReadEntity(serviceLink)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    err = dao.AddServiceLink(*serviceLink)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    response.WriteHeader(http.StatusCreated)
}

func (handler *serviceHandler)findServiceLink(request *restful.Request, response *restful.Response){
    serviceName := request.PathParameter("service-name")
    linkType := request.QueryParameter("linkType")
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
        return
    }
    if linkType == ""{
        linkType = "both"
    }

    serviceAndServiceLink, err := dao.GetServiceAndServiceLink(serviceName, strings.ToLower(linkType))
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(serviceAndServiceLink)
}

func (handler *serviceHandler)findServiceTree(request *restful.Request, response *restful.Response){
    serviceName := request.PathParameter("service-name")
    linkType := request.QueryParameter("linkType")
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
        return
    }
    if linkType == ""{
        linkType = "target"
    }

    serviceTree, err := dao.GetServiceTree(serviceName, strings.ToLower(linkType))
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(serviceTree)
}