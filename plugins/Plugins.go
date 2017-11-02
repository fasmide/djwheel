package plugins

import (
	"io"
	"log"
)

type Plugin interface {
	Priority() int
	Write(io.Writer)
}

var plugins map[string]Plugin

func init() {
	plugins = make(map[string]Plugin)
}

func RegisterPlugin(n string, p Plugin) {
	log.Printf("Hello plugin %s", n)
	plugins[n] = p
}

func Write(to io.Writer) {
	// TODO We need a mixer
	// TODO we should figure out what plugin is active
	// for now, just spectrum

	plugins[priorityPlugin()].Write(to)
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
