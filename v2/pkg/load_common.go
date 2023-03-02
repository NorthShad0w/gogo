package pkg

import (
	"encoding/json"
	"github.com/chainreactors/ipcs"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/parsers/iutils"
	"regexp"
	"strings"
)

var (
	NameMap = ipcs.NameMap
	PortMap = ipcs.PortMap
	TagMap  = ipcs.TagMap
	//WorkFlowMap    map[string][]*Workflow
	Extractor      []*parsers.Extractor
	Extractors                                 = make(parsers.Extractors)
	ExtractRegexps map[string][]*regexp.Regexp = map[string][]*regexp.Regexp{}
)

type PortFinger struct {
	Name  string   `json:"name"`
	Ports []string `json:"ports"`
	Tags  []string `json:"tags"`
}

func LoadPortConfig() {
	var portfingers []PortFinger
	err := json.Unmarshal(LoadConfig("port"), &portfingers)

	if err != nil {
		iutils.Fatal("port config load FAIL!, " + err.Error())
	}
	for _, v := range portfingers {
		v.Ports = ipcs.ParsePorts(v.Ports)
		ipcs.NameMap.Append(v.Name, v.Ports...)
		for _, t := range v.Tags {
			ipcs.TagMap.Append(t, v.Ports...)
		}
		for _, p := range v.Ports {
			ipcs.PortMap.Append(p, v.Name)
		}
	}
}

func LoadExtractor() {
	err := json.Unmarshal(LoadConfig("extract"), &Extractor)
	if err != nil {
		iutils.Fatal("extract config load FAIL!, " + err.Error())
	}

	for _, extract := range Extractor {
		extract.Compile()

		ExtractRegexps[extract.Name] = extract.CompiledRegexps
		for _, tag := range extract.Tags {
			if _, ok := ExtractRegexps[tag]; !ok {
				ExtractRegexps[tag] = extract.CompiledRegexps
			} else {
				ExtractRegexps[tag] = append(ExtractRegexps[tag], extract.CompiledRegexps...)
			}
		}
	}
}

func LoadWorkFlow() WorkflowMap {
	var workflows []*Workflow
	var err error
	err = json.Unmarshal(LoadConfig("workflow"), &workflows)
	if err != nil {
		iutils.Fatal("workflow load FAIL, " + err.Error())
	}

	// 设置默认参数
	for _, w := range workflows {
		// 参数默认值
		if w.IpProbe == "" {
			w.IpProbe = Default
		}
		if w.SmartProbe == "" {
			w.SmartProbe = Default
		}
		if w.Ports == "" {
			w.Ports = "top1"
		}
		if w.Mod == "" {
			w.Mod = Default
		}
		if w.File == "" {
			w.File = "auto"
		}
		//if w.Path == "" {
		//	w.Path = "."
		//}
		if w.Exploit == "" {
			w.Exploit = "none"
		}
	}

	var tmpmap = make(map[string][]*Workflow)
	for _, workflow := range workflows {
		tmpmap[strings.ToLower(workflow.Name)] = append(tmpmap[strings.ToLower(workflow.Name)], workflow)
		for _, tag := range workflow.Tags {
			tmpmap[strings.ToLower(tag)] = append(tmpmap[strings.ToLower(tag)], workflow)
		}
	}
	return tmpmap
}

type WorkflowMap map[string][]*Workflow

func (m WorkflowMap) Choice(name string) []*Workflow {
	var workflows []*Workflow
	name = strings.TrimSpace(name)
	names := strings.Split(name, ",")
	for _, n := range names {
		workflows = append(workflows, m[strings.ToLower(n)]...)
	}
	return workflows
}
