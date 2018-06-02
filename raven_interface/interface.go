// Nightswatch agent plugins interface
package raven_interface

type ReportType int
type ReportValueType int

const (
	ReportGroup ReportType = iota
	ReportSingle
)

const (
	ReportValStr ReportValueType = iota
	ReportValFloat
	ReportValInt
)

type Report struct {
	Key       string
	ValueType ReportValueType
	Value     interface{}
}

type PluginReport struct {
	PluginName string
	ReportType ReportType
	// either *Report{} or []*Report
	Report interface{}
}

type Plugin struct {
	Name  string
	Ch    chan *PluginReport
	ErrCh chan *error
}

func NewPlugin(name string, ch chan *PluginReport, errCh chan *error) *Plugin {
	return &Plugin{name, ch, errCh}
}

func (p *Plugin) SingleReport(valueType ReportValueType,
	k string, v interface{}) {
	p.Ch <- &PluginReport{
		p.Name,
		ReportSingle,
		&Report{k, valueType, v}}
}

// TODO: now the reports param must be map[string]interface{}
// not map[string]float or map[string]string
// or any value the user wants
func (p *Plugin) GroupReport(valueType ReportValueType,
	reports map[string]interface{}) {
	s := make([]*Report, len(reports), len(reports))
	i := 0
	for k, v := range reports {
		s[i] = &Report{k, valueType, v}
		i++
	}
	p.Ch <- &PluginReport{
		p.Name,
		ReportGroup,
		s}
}

func (p *Plugin) ReportError(err error) {
	p.ErrCh <- &err
}
