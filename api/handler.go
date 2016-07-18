package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
    "net/http"
    "strings"
    "fmt"
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

func (handler *serviceHandler)findAllServiceNames(request *restful.Request, response *restful.Response){
    serviceList, err := dao.GetAllService()
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    serviceNameList := make([]string, 0)
    for _, service := range serviceList {
        serviceNameList = append(serviceNameList, service.Name)
    }
    response.WriteEntity(serviceNameList)
}

func (handler *serviceHandler)findServiceByGitIds(request *restful.Request, response *restful.Response){
    gitIdList := request.QueryParameter("gitIdList")
    if gitIdList == ""{
        response.WriteErrorString(http.StatusBadRequest, "git Id列表不能为空")
        return
    }
    gitIds := strings.Split(gitIdList, "|")
    serviceList, err := dao.GetServiceListByGitId(gitIds)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(serviceList)
}

func (handler *serviceHandler)createService(request *restful.Request, response *restful.Response){
    service := &dao.Service{Name: request.PathParameter("service-name")}
    err := request.ReadEntity(service)
    if err != nil{
        fmt.Printf("try create service but return error:%v", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    err = dao.AddService(*service, "replace")
    if err != nil{
        fmt.Printf("try create service but return error:%v", err)
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

func (handler *serviceHandler)findServiceListLink(request *restful.Request, response *restful.Response){
    serviceNameList := strings.Split(request.QueryParameter("serviceList"), "|")
    if len(serviceNameList) == 0 {
        response.WriteErrorString(http.StatusBadRequest, "服务名称列表不能为空")
        return
    }
    // serviceList
    serviceList, err := dao.GetServiceList(serviceNameList)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    // serviceLinkList
    serviceLinkList, err := dao.GetServiceListEachLink(serviceNameList)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    serviceAndServiceLink := &dao.ServiceAndServiceLink{serviceList, serviceLinkList}

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

func (handler *serviceHandler)findServiceListTree(request *restful.Request, response *restful.Response){
    serviceNames := request.QueryParameter("serviceList")
    linkType := strings.ToLower(request.QueryParameter("linkType"))
    serviceType := strings.ToLower(request.QueryParameter("serviceType"))
    serviceCategory := strings.ToLower(request.QueryParameter("serviceCategory"))
    if serviceNames == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称列表不能为空")
        return
    }
    if linkType == "" {
        linkType = "target"
    }
    if serviceType == "" {
        serviceType = "all"
    }
    if serviceCategory == "" {
        serviceCategory = "all"
    }

    serviceListTree, err := dao.GetServiceListTree(strings.Split(serviceNames, "|"), linkType, serviceType, serviceCategory)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    response.WriteEntity(serviceListTree)
}