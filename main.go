package main

import (
	"log"
	p "port-mapping/parse"
	s "port-mapping/service"
	l "port-mapping/wathcer"
)

var (
	application = &MappingApplication{}
)

type MappingApplication struct {
	rules p.MappingRules

	server *s.MappingServer
}

func main() {

	application.startUp()

	l.WatchFileChanged(onWrite())

	<-make(chan any)

}

func (app *MappingApplication) startUp() {
	rules, err := p.GetRules()
	if err != nil {
		log.Fatal("Enable failed with error toml file:", err)
	}
	app.rules = rules
	app.server = s.CreateMappingServer()
	app.server.RegisterTasks(app.rules)
}

func (app *MappingApplication) refresh(newRules p.MappingRules) {
	oldRules := app.rules
	app.rules = newRules
	added, removed := diffRules(oldRules, newRules)
	log.Println("Add rules:", added, "Remove rules:", removed)
	app.server.UnRegisterTasks(removed)
	app.server.RegisterTasks(added)
}

func diffRules(oldRules p.MappingRules, newRules p.MappingRules) (p.MappingRules, p.MappingRules) {
	var added []p.ForwardingPortRule
	var removed []p.ForwardingPortRule
	for _, rule := range newRules {
		if !oldRules.Contains(rule) {
			added = append(added, rule)
		}
	}

	for _, rule := range oldRules {
		if !newRules.Contains(rule) {
			removed = append(removed, rule)
		}
	}
	return added, removed
}

func onWrite() func() {
	return func() {
		rules, err := p.GetRules()
		if err != nil {
			log.Println("Flush group file failed with error toml file:", err)
		}
		application.refresh(rules)
	}
}
