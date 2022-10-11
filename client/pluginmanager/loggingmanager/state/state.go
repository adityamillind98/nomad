package state

import pstructs "github.com/hashicorp/nomad/plugins/shared/structs"

// PluginState is used to store the logging plugin manager's state across
// restarts of the agent
type PluginState struct {
	// ReattachConfigs are the set of reattach configs for plugins launched by
	// the logging plugin manager
	ReattachConfigs map[string]*pstructs.ReattachConfig
}
