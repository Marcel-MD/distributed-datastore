package models

type Instance struct {
	IsLeader bool   `json:"is_leader"`
	Host     string `json:"host"`
	HttpPort string `json:"http_port"`
	TcpPort  string `json:"tcp_port"`
}
