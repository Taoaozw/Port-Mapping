package parse

import (
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"log"
)

const (
	groupFilePath  = "config/group.toml"
	activeFilePath = "config/mapping.toml"
)

type MappingRules []ForwardingPortRule

type ForwardingPortRule struct {
	Name       string
	LocalPort  int
	RemotePort int
	RemoteHost string
}

type RuleGroups struct {
	GroupRules []*GroupRules
}

type GroupRules struct {
	Name       string
	RemoteHost string
	Rules      MappingRules
}

type activeMappingPorts struct {
	ActiveGroups []string
	MappingRules MappingRules
}

func GetRules() (MappingRules, error) {
	activeRules, err := activeMappingRules(activeFilePath)
	if err != nil {
		return nil, err
	}
	groups, err := mappingRuleGroups(groupFilePath)
	if err != nil {
		return nil, err
	}
	activeGroupMappings(&activeRules, &groups)
	return activeRules.MappingRules, err
}

func (rule ForwardingPortRule) isSameWith(other ForwardingPortRule) bool {
	return rule.Name == other.Name && rule.LocalPort == other.LocalPort && rule.RemotePort == other.RemotePort && rule.RemoteHost == other.RemoteHost
}

func (rules MappingRules) Contains(rule ForwardingPortRule) bool {
	for _, r := range rules {
		if r.isSameWith(rule) {
			return true
		}
	}
	return false
}

func activeGroupMappings(activeRules *activeMappingPorts, ruleGroups *RuleGroups) {
	if len(activeRules.ActiveGroups) > 0 {
		for _, groupName := range activeRules.ActiveGroups {
			for _, rule := range ruleGroups.GroupRules {
				if rule.Name == groupName {
					activeRules.MappingRules = append(activeRules.MappingRules, rule.Rules...)
				}
			}
		}
	}
}

func mappingRuleGroups(path string) (RuleGroups, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ruleGroups RuleGroups
	err = toml.Unmarshal(file, &ruleGroups)
	if err != nil {
		log.Println("Read toml file failed:", err)
	}
	log.Println("load group service port :", ruleGroups.GroupRules)
	ruleGroups.fillData()
	return ruleGroups, err
}

func activeMappingRules(path string) (activeMappingPorts, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var rules activeMappingPorts
	err = toml.Unmarshal(file, &rules)
	if err != nil {
		log.Println("Read toml file failed:", err)
	}
	log.Println("load active service port :", rules.MappingRules)
	return rules, err
}

func (rg *RuleGroups) fillData() {
	for _, rules := range rg.GroupRules {
		for i := range rules.Rules {
			rules.Rules[i].RemoteHost = rules.RemoteHost
		}
	}
}
