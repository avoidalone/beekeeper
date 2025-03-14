package config

import (
	"reflect"

	"github.com/ethersphere/beekeeper/pkg/orchestration"
)

// Cluster represents cluster configuration
type Cluster struct {
	// parent to inherit settings from
	*Inherit `yaml:",inline"`
	// Cluster configuration
	Name               *string                      `yaml:"name,omitempty"`
	Namespace          *string                      `yaml:"namespace"`
	DisableNamespace   *bool                        `yaml:"disable-namespace"`
	UseStaticEndpoints *bool                        `yaml:"use-static-endpoints"`
	APIDomain          *string                      `yaml:"api-domain"`
	APIDomainInternal  *string                      `yaml:"api-domain-internal"`
	APIInsecureTLS     *bool                        `yaml:"api-insecure-tls"`
	APIScheme          *string                      `yaml:"api-scheme"`
	Funding            *Funding                     `yaml:"funding"`
	NodeGroups         *map[string]ClusterNodeGroup `yaml:"node-groups"`
}

func (b Cluster) GetParentName() string {
	if b.Inherit != nil {
		return b.Inherit.ParentName
	}
	return ""
}

// ClusterNodeGroup represents node group in the cluster
type ClusterNodeGroup struct {
	cluster       *Cluster
	Mode          string         `yaml:"mode"`
	BeeConfig     string         `yaml:"bee-config"`
	Config        string         `yaml:"config"`
	Count         int            `yaml:"count"`
	Nodes         []ClusterNode  `yaml:"nodes"`
	NodeEndpoints []NodeEndpoint `yaml:"endpoints"`
}

// ClusterNode represents node in the cluster
type ClusterNode struct {
	Name      string                      `yaml:"name"`
	Bootnodes string                      `yaml:"bootnodes"`
	LibP2PKey string                      `yaml:"libp2p-key"`
	SwarmKey  *orchestration.EncryptedKey `yaml:"swarm-key"`
}

type NodeEndpoint struct {
	Name   string `yaml:"name"`
	APIURL string `yaml:"api-url"`
}

// Export exports Cluster to orchestration.ClusterOptions, skipping all other extra fields
func (c *Cluster) Export() (o orchestration.ClusterOptions) {
	localVal := reflect.ValueOf(c).Elem()
	localType := reflect.TypeOf(c).Elem()
	remoteVal := reflect.ValueOf(&o).Elem()

	for i := 0; i < localVal.NumField(); i++ {
		localField := localVal.Field(i)
		if localField.IsValid() && !localField.IsNil() {
			localFieldVal := localVal.Field(i).Elem()
			localFieldName := localType.Field(i).Name

			remoteFieldVal := remoteVal.FieldByName(localFieldName)
			if remoteFieldVal.IsValid() && remoteFieldVal.Type() == localFieldVal.Type() {
				remoteFieldVal.Set(localFieldVal)
			}
		}
	}

	return remoteVal.Interface().(orchestration.ClusterOptions)
}

// GetName returns cluster name
func (c *Cluster) GetName() string {
	if c.Name == nil {
		return "noname"
	}
	return *c.Name
}

// GetNamespace returns cluster namespace
func (c *Cluster) GetNamespace() string {
	if c.Namespace == nil {
		return "nonamespace"
	}
	return *c.Namespace
}

// GetNodeGroups returns cluster node groups
func (c *Cluster) GetNodeGroups() map[string]ClusterNodeGroup {
	if c.NodeGroups == nil {
		return nil
	}

	nodeGroups := *c.NodeGroups
	for key, group := range nodeGroups {
		group.cluster = c // Set the reference to the parent cluster
		nodeGroups[key] = group
	}

	return nodeGroups
}

// IsUsingStaticEndpoints
func (c *Cluster) IsUsingStaticEndpoints() bool {
	if c.UseStaticEndpoints == nil {
		return false
	}
	return *c.UseStaticEndpoints
}

func (ng *ClusterNodeGroup) GetEndpoints() map[string]NodeEndpoint {
	endpoints := make(map[string]NodeEndpoint)
	for _, endpoint := range ng.NodeEndpoints {
		endpoints[endpoint.Name] = endpoint
	}
	return endpoints
}
