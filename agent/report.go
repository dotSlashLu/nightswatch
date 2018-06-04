package main

import (
	"fmt"
	ri "github.com/dotSlashLu/nightswatch/raven_interface"
	"strconv"
)

func processReport(r *ri.Report) {
	switch r.ValueType {
	case ri.ReportValStr:
		fmt.Printf("report str k: %v, v: %v\n", r.Key, r.Value)
		q.Push(r.Key, r.Value.(string))
	case ri.ReportValFloat:
		fmt.Printf("report float k: %v, v: %v\n", r.Key, r.Value)
		q.Push(r.Key, strconv.FormatFloat(r.Value.(float64), 'f', 6, 64))
	case ri.ReportValInt:
		fmt.Printf("report int k: %v, v: %v\n", r.Key, r.Value)
		q.Push(r.Key, string(r.Value.(int)))
	}

}

func processReportRaw(rs *ri.PluginReport) {
	switch rs.ReportType {
	case ri.ReportGroup:
		reports := rs.Report.([]*ri.Report)
		for _, r := range reports {
			processReport(r)
		}
	case ri.ReportSingle:
		report := rs.Report.(*ri.Report)
		processReport(report)
	default:
		fmt.Printf("plugin %s reported with an unrecognized type",
			rs.PluginName)
	}
}
