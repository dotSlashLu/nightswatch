package common

func LoadPlugin(pluginName string, ch chan *ri.PluginReport,
	errChan chan *error) {
	path := "/root/go/src/github.com/dotSlashLu/nightswatch/raven/" +
		pluginName + ".so"
	so, err := plugin.Open(path)
	if err != nil {
		panic(fmt.Sprintf("error loading plugin %s at %s: %v",
			pluginName, path, err))
	}
	Report, err := so.Lookup("Report")
	report, ok := Report.(func(chan *ri.PluginReport, chan *error))
	if !ok {
		fmt.Printf("Plugin %s not conforming to interface, ignoring\n", pluginName)
		return
	}
	go report(ch, errChan)
}
