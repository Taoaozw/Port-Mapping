package service

import (
	"fmt"
	"log"
	"net"
	n "port-mapping/network"
	p "port-mapping/parse"
	"sync"
	"time"
)

type MappingTask struct {
	Rule          p.ForwardingPortRule
	ClientMapLock sync.Mutex
	PortListener  net.Listener
	ClientMap     map[string]*MappingClient
}

func (_self *MappingTask) StartJob() {
	_self.ClientMap = make(map[string]*MappingClient, 200)
	sourceAddr := fmt.Sprint("127.0.0.1", ":", _self.Rule.LocalPort)
	destAddr := fmt.Sprint(_self.Rule.RemoteHost, ":", _self.Rule.RemotePort)

	var err error
	_self.PortListener, err = n.NewTCP(sourceAddr)

	if err != nil {
		log.Println("启动监听 ", sourceAddr, " 出错：", err)
	}

	log.Println("启动端口转发，从 ", sourceAddr, " 到 ", destAddr)

	_self.doTcpForward(destAddr)

}

func (_self *MappingTask) doTcpForward(destAddr string) {

	for {
		fromConnection, err := _self.PortListener.Accept()
		if err != nil {
			log.Println("Forward Accept err:", err.Error())
			_self.StopJob()
			break
		}

		toConnection, err := net.DialTimeout("tcp", destAddr, 30*time.Second)

		if err != nil {
			log.Print("转发出现异常 Forward to Destination Addr err:", err.Error())
			continue
		}
		forwardClient := &MappingClient{fromConnection, toConnection, nil, _self.ClosedCallBack}
		forwardClient.StartForward()
		_self.registerClient(_self.getClientId(fromConnection), forwardClient)
	}
}

func (_self *MappingTask) ClosedCallBack(srcConn net.Conn, destConn net.Conn) {

	_self.UnRegistryClient(_self.getClientId(srcConn))
}

func (_self *MappingTask) UnRegistryClient(srcAddr string) {
	_self.ClientMapLock.Lock()
	defer _self.ClientMapLock.Unlock()
	delete(_self.ClientMap, srcAddr)
	log.Println("UnRegistryClient srcAddr: ", srcAddr)
}

func (_self *MappingTask) getClientId(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

func (_self *MappingTask) StopJob() {
	_self.PortListener.Close()
	for srcAddr, client := range _self.ClientMap {
		log.Println("停止真实用户连接：", srcAddr)
		client.StopForward()
	}
	_self.ClientMap = nil
}

func (_self *MappingTask) registerClient(srcAddr string, forwardClient *MappingClient) {
	_self.ClientMapLock.Lock()
	defer _self.ClientMapLock.Unlock()

	_self.ClientMap[srcAddr] = forwardClient

}
