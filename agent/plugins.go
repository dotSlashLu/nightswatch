package main

import (
	"fmt"
	ri "github.com/dotSlashLu/nightswatch/agent/raven_interface"
	"plugin"
)

func loadPlugin(pluginName string, ch chan *ri.PluginReport,
	errChan chan *error) {
	path := cfg.Plugins.Directory +
		pluginName + ".so"
	so, err := plugin.Open(path)
	if err != nil {
		panic(fmt.Sprintf("error loading plugin %s at %s: %v",
			pluginName, path, err))
	}
	Report, err := so.Lookup("Report")
	report, ok := Report.(func(chan *ri.PluginReport, chan *error))
	if !ok {
		fmt.Printf("Plugin %s not conforming to interface, ignoring\n",
			pluginName)
		return
	}
	go report(ch, errChan)
}

func loadPlugins() {
	plugins := cfg.Plugins
	ch := make(chan *ri.PluginReport)
	errCh := make(chan *error)
	for _, plugin := range plugins.Names {
		loadPlugin(plugin, ch, errCh)
	}

	// poll
	for {
		fmt.Println("waiting for report")
		select {
		case reports := <-ch:
			processReportRaw(reports)
		case err := <-errCh:
			fmt.Printf("plugin error: %s\n", (*err).Error())
		}
	}
}
