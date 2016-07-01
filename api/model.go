package api

import "chaos/dao"

// 服务的信息和服务之间的调用关系
type ServiceAndServiceLink struct {
    ServiceList         []dao.Service       `json:"serviceList,omitempty"`
    ServiceLinkList     []dao.ServiceLink   `json:"serviceLinkList,omitempty"`
}