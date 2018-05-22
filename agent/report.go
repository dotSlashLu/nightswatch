package main

import (
	"fmt"
	ri "github.com/dotSlashLu/nightswatch/raven_interface"
)

func processReport(r *ri.Report) {
	switch r.ValueType {
	case ri.ReportValStr:
		fmt.Printf("report str k: %v, v: %v\n", r.Key, r.Value)
	case ri.ReportValFloat:
		fmt.Printf("report float k: %v, v: %v\n", r.Key, r.Value)
	case ri.ReportValInt:
		fmt.Printf("report int k: %v, v: %v\n", r.Key, r.Value)
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
		fmt.Printf("plugin %s reported with an unrecognized type", rs.PluginName)
	}
}
