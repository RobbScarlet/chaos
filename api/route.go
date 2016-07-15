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
    // 获取所有服务名称
    ws.Route(
        ws.GET("/names").To(h.findAllServiceNames).
        Doc("get all service name list").
        Writes([]string{}),
    )
    // 获取通过gitId,获取服务信息
    ws.Route(
        ws.GET("/info").To(h.findServiceByGitIds).
        Doc("get service info by gitids").
        // gitid列表,多个以'|'分隔
        Param(ws.QueryParameter("gitIdList", "git id list, use '|' split").DataType("string")).
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
        // 服务名称
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        // 查询调用关系类型:source为查询父, target为查询子,both为查询父和子,默认是both
        Param(ws.QueryParameter("linkType", "service link type('both', 'source', 'target'),default is 'both'").DataType("string")).
        Writes(dao.ServiceAndServiceLink{}),
    )
    // 获取多个服务相互之间的直接调用关系
    ws.Route(
        ws.GET("/link/eachother").To(h.findServiceListLink).
        Doc("get service list info and links").
        // 服务名称列表,多个以'|'分隔
        Param(ws.QueryParameter("serviceList", "service name list use '|' split").DataType("string")).
        Writes(dao.ServiceAndServiceLink{}),
    )
    // 获取服务的调用关系树
    ws.Route(
        ws.GET("/tree/service/{service-name}").To(h.findServiceTree).
        Doc("get service and tree by service name").
        // 服务名称
        Param(ws.PathParameter("service-name", "name of the service").DataType("string")).
        // 查询关系类型:source为递归查询父, target为递归查询子,默认target
        Param(ws.QueryParameter("linkType", "service link type('source', 'target'), default is 'target'").DataType("string")).
        Writes(dao.ServiceTree{}),
    )
    // 获取多个服务的调用关系树
    ws.Route(
        ws.GET("/tree/services").To(h.findServiceListTree).
        Doc("get service list info and service tree only contain which type").
        // 服务名称列表,多个以'|'分隔
        Param(ws.QueryParameter("serviceList", "service name list use '|' split").DataType("string")).
        // 查询关系类型:source为递归查询父, target为递归查询子,默认target
        Param(ws.QueryParameter("linkType", "service link type('source', 'target'), default is 'target'").DataType("string")).
        // 查询保留的递归服务的类型,'http', 'thrift'等,多个以'|'分隔
        Param(ws.QueryParameter("serviceType", "default is all, the value could be 'http', 'thrift' and so on, more use '|' split").DataType("string")).
        // 查询保留的递归服务的种类,'base'和'application',不传的话则为保留所有
        Param(ws.QueryParameter("serviceCategory", "default is all, the value could be 'base' or 'application'").DataType("string")).
        Writes(dao.ServiceListTree{}),
    )

    container.Add(ws)
    return
}
