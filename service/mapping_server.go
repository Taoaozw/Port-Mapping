package service

import (
	"log"
	"sync"
)
import p "port-mapping/parse"

type MappingServer struct {
	Lock    sync.Mutex
	TaskMap map[string]*MappingTask
}

func CreateMappingServer() *MappingServer {
	return &MappingServer{
		TaskMap: make(map[string]*MappingTask, 200),
	}
}

func (_self *MappingServer) RegisterTasks(rules p.MappingRules) {
	for _, rule := range rules {
		forWardJob := new(MappingTask)
		forWardJob.Rule = rule
		go forWardJob.StartJob()
		_self.registryJob(rule, forWardJob)
	}
}

func (_self *MappingServer) UnRegisterTasks(rules p.MappingRules) {
	_self.Lock.Lock()
	defer _self.Lock.Unlock()
	for _, rule := range rules {
		if v, ok := _self.TaskMap[rule.Name]; ok {
			v.StopJob()
			delete(_self.TaskMap, rules[0].Name)
		} else {
			log.Printf("UnRegisterTasks failed with rule:%v ,beacuse MappingServer dont find this task.\n Current Tasks [%v]", rule, _self.TaskMap)
		}
	}
}

func (_self *MappingServer) registryJob(rule p.ForwardingPortRule, forWardJob *MappingTask) {
	_self.Lock.Lock()
	defer _self.Lock.Unlock()
	_self.TaskMap[rule.Name] = forWardJob

}
