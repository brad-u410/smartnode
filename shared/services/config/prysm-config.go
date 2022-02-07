package config

import (
	"fmt"
	"runtime"
)

// v2.0.6
const prysmBnTagAmd64 string = "prysmaticlabs/prysm-beacon-chain:HEAD-e26cde-debug"
const prysmVcTagAmd64 string = "prysmaticlabs/prysm-validator:HEAD-e26cde-debug"
const prysmTagArm64 string = "rocketpool/prysm:v2.0.6"
const defaultPrysmRpcPort uint16 = 5053
const defaultPrysmOpenRpcPort bool = false

// Configuration for Prysm
type PrysmConfig struct {
	// Common parameters shared across clients
	CommonParams *ConsensusCommonParams

	// Common parameters that Prysm doesn't support and should be hidden
	UnsupportedCommonParams []string

	// The RPC port for BN / VC connections
	RpcPort *Parameter

	// Toggle for forwarding the RPC API outside of Docker
	OpenRpcPort *Parameter

	// The Docker Hub tag for the Prysm BN
	BnContainerTag *Parameter

	// The Docker Hub tag for the Prysm VC
	VcContainerTag *Parameter

	// Custom command line flags for the BN
	AdditionalBnFlags *Parameter

	// Custom command line flags for the VC
	AdditionalVcFlags *Parameter
}

// Generates a new Prysm configuration
func NewPrysmConfig(commonParams *ConsensusCommonParams) *PrysmConfig {
	return &PrysmConfig{
		CommonParams: commonParams,

		UnsupportedCommonParams: []string{
			checkpointSyncUrlID,
		},

		RpcPort: &Parameter{
			ID:                   "rpcPort",
			Name:                 "RPC Port",
			Description:          "The port Prysm should run its JSON-RPC API on.",
			Type:                 ParameterType_Uint16,
			Default:              defaultPrysmRpcPort,
			AffectsContainers:    []ContainerID{ContainerID_Eth2, ContainerID_Validator},
			EnvironmentVariables: []string{"BN_RPC_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		OpenRpcPort: &Parameter{
			ID:                   "openRpcPort",
			Name:                 "Open RPC Port",
			Description:          "Enable this to open Prysm's API ports to your local network, so other machines can access it too.",
			Type:                 ParameterType_Bool,
			Default:              defaultPrysmOpenRpcPort,
			AffectsContainers:    []ContainerID{ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_OPEN_RPC_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		BnContainerTag: &Parameter{
			ID:                   "bnContainerTag",
			Name:                 "Beacon Node Container Tag",
			Description:          "The tag name of the Prysm Beacon Node container you want to use on Docker Hub.",
			Type:                 ParameterType_String,
			Default:              getPrysmBnTag(),
			AffectsContainers:    []ContainerID{ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		VcContainerTag: &Parameter{
			ID:                   "vcContainerTag",
			Name:                 "Validator Client Container Tag",
			Description:          "The tag name of the Prysm Validator Client container you want to use on Docker Hub.",
			Type:                 ParameterType_String,
			Default:              getPrysmVcTag(),
			AffectsContainers:    []ContainerID{ContainerID_Validator},
			EnvironmentVariables: []string{"VC_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		AdditionalBnFlags: &Parameter{
			ID:                   "additionalBnFlags",
			Name:                 "Additional Beacon Node Flags",
			Description:          "Additional custom command line flags you want to pass Prysm's Beacon Node, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 ParameterType_String,
			Default:              "",
			AffectsContainers:    []ContainerID{ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},

		AdditionalVcFlags: &Parameter{
			ID:                   "additionalVcFlags",
			Name:                 "Additional Validator Client Flags",
			Description:          "Additional custom command line flags you want to pass Prysm's Validator Client, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 ParameterType_String,
			Default:              "",
			AffectsContainers:    []ContainerID{ContainerID_Validator},
			EnvironmentVariables: []string{"VC_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the container tag for the Prysm BN based on the current architecture
func getPrysmBnTag() string {
	if runtime.GOARCH == "arm64" {
		return prysmTagArm64
	} else if runtime.GOARCH == "amd64" {
		return prysmBnTagAmd64
	} else {
		panic(fmt.Sprint("Prysm doesn't support architecture %s", runtime.GOARCH))
	}
}

// Get the container tag for the Prysm VC based on the current architecture
func getPrysmVcTag() string {
	if runtime.GOARCH == "arm64" {
		return prysmTagArm64
	} else if runtime.GOARCH == "amd64" {
		return prysmVcTagAmd64
	} else {
		panic(fmt.Sprint("Prysm doesn't support architecture %s", runtime.GOARCH))
	}
}