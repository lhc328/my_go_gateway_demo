package dao

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"loadbalance" description:"负载均衡策略"`
	AccessControl *AccessControl `json:"accesscontrol" description:"权限控制"`
}
