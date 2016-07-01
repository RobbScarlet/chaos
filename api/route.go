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
        ws.PUT("/base").To(h.createService).
        Doc("create service").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Reads(dao.Service{}),
    )
    // 获取一个服务
    ws.Route(
        ws.GET("/base/{service-name}").To(h.findService).
        Doc("get service by service name").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Writes(dao.Service{}),
    )
    // 创建一个服务调用关系
    ws.Route(
        ws.PUT("/link").To(h.createServiceLink).
        Doc("create service link").
        Reads(dao.ServiceLink{}),
    )
    // 获取一个服务调用关系
    ws.Route(
        ws.GET("/link/{service-name}").To(h.findServiceLink).
        Doc("get service by service name").
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        Param(ws.QueryParameter("linkType", "service link type('both', 'source', 'target'),default is 'both'").DataType("string")).
        Writes(ServiceAndServiceLink{}),
    )

    container.Add(ws)
    return
}
