package utils

import (
	"plugin"

	"github.com/rs/zerolog/log"

	"github.com/amazingchow/photon-dance-mapreduce/backend/storage"
)

var (
	// Plugin is the plugin file
	Plugin = "None"
)

// LoadPlugin loads the MapReduce application from a plugin file.
func LoadPlugin() (func(string, string) []storage.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(Plugin)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot load plugin %s", Plugin)
	}

	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot find Map in plugin %s", Plugin)
	}
	mapf := xmapf.(func(string, string) []storage.KeyValue)

	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot find Reduce in plugin %s", Plugin)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
