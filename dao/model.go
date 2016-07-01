package dao

//
//  服务的基本信息
//
type Service struct {
    // 服务名称(唯一性)
    Name            string              `json:"name,omitempty"`
    // 服务全名
    FullName        string              `json:"fullName,omitempty"`
    // 服务别名
    Alias           string              `json:"alias,omitempty"`
    // Docker中的名称
    DockerName      string              `json:"dockerName,omitempty"`
    // 服务所属的命名空间
    Namespace       string              `json:"namespace,omitempty"`
    // 服务类型(Thrift, Mysql, NodeJS, Http等)
    Type            string              `json:"type,omitempty"`
    // 服务种类(基础服务Base,应用服务Application)
    Category        string              `json:"category,omitempty"`
    // 服务的端口
    Ports           []string            `json:"ports,omitempty"`
    // 服务的所有者
    Owner           string              `json:"owner,omitempty"`
    // 服务的所属部门
    Department      string              `json:"department,omitempty"`
    // 服务描述
    Description     string              `json:"description,omitempty"`
    //
    // 以下属性只有Service.Category=Application的服务才会有
    //
    // 编程语言
    Language        string              `json:"language,omitempty"`
    // Git地址
    GitAddress      string              `json:"gitAddress,omitempty"`
    // GitId
    GitId           string              `json:"gitId,omitempty"`
    // 项目代码所在目录
    ProjectDir      string              `json:"projectDir,omitempty"`
    // 如果接入Loom并提供服务
    LoomId          string              `json:"loomId,omitempty"`
}
//
//  服务基本信息列表
//
type ServiceList struct{
    serviceList     []Service           `json:"serviceList,omitempty"`
}
//
//  服务的调用关系信息
//
type ServiceLink struct {
    // 服务调用者名称
    Source          string              `json:"source,omitempty"`
    // 服务被调用者名称
    Target          string              `json:"target,omitempty"`
}
//
//  服务的调用关系列表
//
type ServiceLinkList struct {
    // 服务调用者名称
    Source          string              `json:"source,omitempty"`
    // 服务被调用者名称列表
    TargetList      []string            `json:"targetList,omitempty"`
}
