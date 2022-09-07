package model

// Description of a pod
type PodDescription struct {
	Id               string           `validate:"required" json:"id"`
	Name             string           `validate:"required,max=128" json:"name"`
	Image            string           `validate:"required" json:"image"`
	Tag              string           `validate:"required" json:"tag"`
	PortMaps         []PortMap        `validate:"required" json:"portMaps"`
	EnvironmentMaps  []EnvironmentMap `validate:"required" json:"environmentMaps"`
	MountMaps        []MountMap       `validate:"required" json:"mountMaps"`
	Command          string           `json:"command"`
	ManifestFileUsed string           `validate:"required" json:"manifestFileUsed"`
	Limit            Limit            `validate:"required" json:"limit"`
}

// Required,
// Important when running more than one type of game server that uses the same port.
// Example
// Minectaft Server One  (Exposed: 22565, Binding: 22565)
// Minectaft Server Two  (Exposed: 22566, Binding: 22565)
// Minectaft Server One  (Exposed: 30215, Binding: 22565)
type PortMap struct {
	Exposed  int    `validate:"required" json:"exposed"` // Port that is internet facing.
	Binding  int    `validate:"required" json:"binding"` // Port used by the container.
	Protocol string `validate:"required" json:"protocol"`
}

// Environment key/value for a container
type EnvironmentMap struct {
	Name  string `validate:"required" json:"name"`
	Value string `validate:"required" json:"value"`
}

// Mount key/value for a container
type MountMap struct {
	Source      string `validate:"required" json:"source"`
	Destination string `validate:"required,dir|file" json:"destination"`
}

type Limit struct {
	Memory int `validate:"required,gt=6" json:"memory"`
}
