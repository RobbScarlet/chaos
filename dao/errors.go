package dao

import "errors"

var (
    ServiceNotExist = errors.New("服务不存在")
    ServiceNameIsNil = errors.New("服务名称为空")
    ServiceNameListIsNil = errors.New("服务名称列表为空")
)