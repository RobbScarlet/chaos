package dao

import "errors"

var (
    ServiceNotExist = errors.New("服务不存在")
    ServiceNameIsNil = errors.New("服务名称为空")
    ServiceNameListIsNil = errors.New("服务名称列表为空")
    ServiceGitIdIsNil = errors.New("Git Id为空")
    ServiceGitIdListIsNil = errors.New("Git Id列表为空")
)