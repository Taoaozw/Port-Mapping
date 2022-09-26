package parse

import (
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"port-mapping/util"
)

const (
	groupFilePath  = "config/group.toml"
	activeFilePath = "config/mapping.toml"

	DISABLE RuleState = 0
	ENABLE  RuleState = 1
)

type RuleState int
type MappingRules []ForwardingPortRule

type ForwardingPortRule struct {
	Name       string
	State      RuleState
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
	State      RuleState
}

type activeMappingPorts struct {
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
	composeGroupMappings(&activeRules, &groups)
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

func composeGroupMappings(activeRules *activeMappingPorts, ruleGroups *RuleGroups) {
	for _, rule := range ruleGroups.GroupRules {
		activeRules.MappingRules = append(activeRules.MappingRules, rule.Rules...)
	}
}

func mappingRuleGroups(path string) (RuleGroups, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var ruleGroups RuleGroups
	err = toml.Unmarshal(file, &ruleGroups)
	if err != nil {
		log.Println("Read toml file failed:", err)
	}
	ruleGroups.GroupRules = util.Filter(ruleGroups.GroupRules, func(rule *GroupRules) bool {
		return rule.State == ENABLE
	})
	log.Printf("load active service group :%v", ruleGroups.GroupRules)
	for _, v := range ruleGroups.GroupRules {
		v.Rules = util.Filter(v.Rules, func(rule ForwardingPortRule) bool {
			return rule.State == ENABLE
		})
	}
	log.Println("load group service port :", ruleGroups.GroupRules)
	ruleGroups.fillData()
	return ruleGroups, err
}

func activeMappingRules(path string) (activeMappingPorts, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var rules activeMappingPorts
	err = toml.Unmarshal(file, &rules)
	if err != nil {
		log.Println("Read toml file failed:", err)
	}
	rules.MappingRules = util.Filter[ForwardingPortRule](rules.MappingRules, func(rule ForwardingPortRule) bool {
		return rule.State == ENABLE
	})

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
