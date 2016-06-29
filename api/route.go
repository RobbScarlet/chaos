package api

import (
    "github.com/emicklei/go-restful"
    "chaos/dao"
)

func Register(container *restful.Container) (err error) {
    h := NewServiceHandler()
    ws := new(restful.WebService)
    ws.
        Path("/chaos/services").
        Doc("Manage Service").
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)
    // 创建一个服务
    ws.Route(
        ws.PUT("/{service-name}").To(h.createService).
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

    container.Add(ws)
    return
}
