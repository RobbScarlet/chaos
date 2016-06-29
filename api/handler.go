package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
    "net/http"
    "fmt"
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
    fmt.Println(serviceName)
    if serviceName == ""{
        response.WriteErrorString(http.StatusBadRequest, "服务名称不能为空")
    }
    service, err := dao.GetService(serviceName)
    if err != nil{
        response.WriteError(http.StatusInternalServerError, err)
    }
    response.WriteEntity(&service)
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
