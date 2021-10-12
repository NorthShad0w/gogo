package templates

import (
	"getitle/src/nuclei/protocols"
	"getitle/src/nuclei/protocols/executer"
	"getitle/src/nuclei/protocols/http"
	"strings"
)

type Template struct {
	Id     string `json:"id"`
	Finger string `json:"finger"`
	Info   struct {
		Name      string `json:"name"`
		Author    string `json:"author"`
		Severity  string `json:"severity"`
		Reference string `json:"reference"`
		Vendor    string `json:"vendor"`
		Tags      string `json:"tags"`
	} `json:"info"`
	RequestsHttp []http.Request `json:"requests"`
	//RequestsTCP []tcp.Request `json:"network"`
	// TotalRequests is the total number of requests for the template.
	TotalRequests int `yaml:"-" json:"-"`
	// Executor is the actual template executor for running template requests
	Executor *executer.Executer `yaml:"-" json:"-"`
}

func (t *Template) GetTags() []string {
	if t.Info.Tags != "" {
		return strings.Split(t.Info.Tags, ",")
	}
	return []string{}
}

func (t *Template) Compile() error {
	options := protocols.ExecuterOptions{}
	var requests []protocols.Request
	var err error
	if len(t.RequestsHttp) > 0 {
		for _, req := range t.RequestsHttp {
			requests = append(requests, &req)
		}
		t.Executor = executer.NewExecuter(requests, &options)
	}
	if t.Executor != nil {
		err = t.Executor.Compile()
		if err != nil {
			return err
		}
		t.TotalRequests += t.Executor.Requests()
	}
	return nil
}

func (t *Template) Execute(url string) (*protocols.Result, bool) {
	res, err := t.Executor.Execute(url)
	if err != nil {
		return nil, false
	}
	return res, true
}
