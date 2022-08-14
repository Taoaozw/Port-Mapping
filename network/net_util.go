package network

import "net"

func NewTCP(addr string) (net.Listener, error) {
	tcpSocket, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return tcpSocket, nil

}
