package api

import "chaos/dao"

type ServiceList struct {
    serviceList []dao.Service           `json:"serviceList,omitempty"`
}
