package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
)

func Register(container *restful.Container) (err error) {
    h := NewServiceHandler()
    ws := new(restful.WebService)
    ws.
        Path("/architecture/services").
        Doc("Manage Service").
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)
    // 创建一个服务
    ws.Route(
        ws.POST("/{service-name}").To(h.createService).
        Doc("create service").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Reads(dao.Service{}),
    )
    // 获取一个服务
    ws.Route(
        ws.GET("/{service-name}").To(h.findService).
        Doc("get service by service name").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Writes(dao.Service{}),
    )
    // 获取所有服务
    ws.Route(
        ws.GET("").To(h.findAllServices).
        Doc("get all service").
        Writes([]dao.Service{}),
    )
    // 创建一个服务调用关系
    ws.Route(
        ws.POST("/link").To(h.createServiceLink).
        Doc("create service link").
        Reads(dao.ServiceLink{}),
    )
    // 获取一个服务调用关系
    ws.Route(
        ws.GET("/link/{service-name}").To(h.findServiceLink).
        Doc("get service and link by service name").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Param(ws.QueryParameter("linkType", "service link type('both', 'source', 'target'),default is 'both'").DataType("string")).
        Writes(dao.ServiceAndServiceLink{}),
    )
    // 获取服务的调用关系树
    ws.Route(
        ws.GET("/tree/{service-name}").To(h.findServiceTree).
        Doc("get service and tree by service name").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Param(ws.QueryParameter("linkType", "service link type('source', 'target'), default is 'target'").DataType("string")).
        Writes(dao.ServiceTree{}),
    )

    container.Add(ws)
    return
}
