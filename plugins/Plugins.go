package plugins

import (
	"io"
	"log"
)

type Plugin interface {
	Priority() int
	WriteTo(io.Writer)
	WheelEvent(int)
}

// plugins holds all registered plugins
var plugins map[string]Plugin

func init() {
	plugins = make(map[string]Plugin)
}

func RegisterPlugin(n string, p Plugin) {
	log.Printf("Hello plugin %s", n)
	plugins[n] = p
}

func WriteTo(to io.Writer) {
	// TODO We need a mixer
	// TODO we should figure out what plugin is active
	// for now, just spectrum

	plugins[priorityPlugin()].WriteTo(to)
}

func HandleWheelEvents(e chan int) {
	for e := range e {
		for _, p := range plugins {
			p.WheelEvent(e)
		}
	}
}

func priorityPlugin() string {
	var highest string
	var priority int
	var currentPriority int

	for name, p := range plugins {
		currentPriority = p.Priority()
		if currentPriority > priority {
			priority = currentPriority
			highest = name
		}
	}
	return highest
}
