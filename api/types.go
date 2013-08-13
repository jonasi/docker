package api

import (
    "fmt"
    "net/http"
    "io"
    "sort"
    "strings"
    "time"
)

type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type JSONMessage struct {
	Status       string     `json:"status,omitempty"`
	Progress     string     `json:"progress,omitempty"`
	ErrorMessage string     `json:"error,omitempty"` //deprecated
	ID           string     `json:"id,omitempty"`
	Time         int64      `json:"time,omitempty"`
	Error        *JSONError `json:"errorDetail,omitempty"`
}

func (e *JSONError) Error() string {
	return e.Message
}

func NewHTTPRequestError(msg string, res *http.Response) error {
	return &JSONError{
		Message: msg,
		Code:    res.StatusCode,
	}
}

func (jm *JSONMessage) Display(out io.Writer) error {
	if jm.Error != nil {
		if jm.Error.Code == 401 {
			return fmt.Errorf("Authentication is required.")
		}
		return jm.Error
	}
	if jm.Time != 0 {
		fmt.Fprintf(out, "[%s] ", time.Unix(jm.Time, 0))
	}
	if jm.ID != "" {
		fmt.Fprintf(out, "%s: ", jm.ID)
	}
	if jm.Progress != "" {
		fmt.Fprintf(out, "%c[2K", 27)
		fmt.Fprintf(out, "%s %s\r", jm.Status, jm.Progress)
	} else {
		fmt.Fprintf(out, "%s\r\n", jm.Status)
	}
	return nil
}

type Version struct {
	Version   string
	GitCommit string `json:",omitempty"`
	GoVersion string `json:",omitempty"`
}

type Info struct {
	Debug              bool
	Containers         int
	Images             int
	NFd                int    `json:",omitempty"`
	NGoroutines        int    `json:",omitempty"`
	MemoryLimit        bool   `json:",omitempty"`
	SwapLimit          bool   `json:",omitempty"`
	IPv4Forwarding     bool   `json:",omitempty"`
	LXCVersion         string `json:",omitempty"`
	NEventsListener    int    `json:",omitempty"`
	KernelVersion      string `json:",omitempty"`
	IndexServerAddress string `json:",omitempty"`
}

type ID struct {
	ID string `json:"Id"`
}

type Copy struct {
	Resource string
}

type Containers struct {
	ID         string `json:"Id"`
	Image      string
	Command    string
	Created    int64
	Status     string
	Ports      string
	SizeRw     int64
	SizeRootFs int64
}

type Auth struct {
	Status string
}

type Search struct {
	Name        string
	Description string
}

type Images struct {
	Repository  string `json:",omitempty"`
	Tag         string `json:",omitempty"`
	ID          string `json:"Id"`
	Created     int64
	Size        int64
	VirtualSize int64
}

type Rmi struct {
	Deleted  string `json:",omitempty"`
	Untagged string `json:",omitempty"`
}

type Wait struct {
	StatusCode int
}

type Run struct {
	ID       string   `json:"Id"`
	Warnings []string `json:",omitempty"`
}

type Top struct {
	Titles    []string
	Processes [][]string
}

type History struct {
	ID        string   `json:"Id"`
	Tags      []string `json:",omitempty"`
	Created   int64
	CreatedBy string `json:",omitempty"`
}

type ImageConfig struct {
	ID string `json:"Id"`
	*Config
}

type Config struct {
	Hostname        string
	User            string
	Memory          int64 // Memory limit (in bytes)
	MemorySwap      int64 // Total memory usage (memory + swap); set `-1' to disable swap
	CpuShares       int64 // CPU shares (relative weight vs. other containers)
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	Tty             bool // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool // Open stdin
	StdinOnce       bool // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Cmd             []string
	Dns             []string
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	VolumesFrom     string
	Entrypoint      []string
	NetworkDisabled bool
}

type Image struct {
	ID              string    `json:"id"`
	Parent          string    `json:"parent,omitempty"`
	Comment         string    `json:"comment,omitempty"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container,omitempty"`
	ContainerConfig Config    `json:"container_config,omitempty"`
	DockerVersion   string    `json:"docker_version,omitempty"`
	Author          string    `json:"author,omitempty"`
	Config          *Config   `json:"config,omitempty"`
	Architecture    string    `json:"architecture,omitempty"`
	Size            int64
}

type ChangeType int

const (
	ChangeModify = iota
	ChangeAdd
	ChangeDelete
)

type Change struct {
	Path string
	Kind ChangeType
}

func (change *Change) String() string {
	var kind string
	switch change.Kind {
	case ChangeModify:
		kind = "C"
	case ChangeAdd:
		kind = "A"
	case ChangeDelete:
		kind = "D"
	}
	return fmt.Sprintf("%s %s", kind, change.Path)
}

type Container struct {
	ID string
	Created time.Time
	Path string
	Args []string
	Config *Config
	State  State
	Image  string
	NetworkSettings *NetworkSettings
	SysInitPath    string
	ResolvConfPath string

	Volumes  map[string]string
	// Store rw/ro in a separate structure to preserve reverse-compatibility on-disk.
	// Easier than migrating older container configs :)
	VolumesRW map[string]bool
}

type State struct {
	Running   bool
	Pid       int
	ExitCode  int
	StartedAt time.Time
	Ghost     bool
}

type NetworkSettings struct {
	IPAddress   string
	IPPrefixLen int
	Gateway     string
	Bridge      string
	PortMapping map[string]PortMapping
}

// String returns a human-readable description of the port mapping defined in the settings
func (settings *NetworkSettings) PortMappingHuman() string {
	var mapping []string
	for private, public := range settings.PortMapping["Tcp"] {
		mapping = append(mapping, fmt.Sprintf("%s->%s", public, private))
	}
	for private, public := range settings.PortMapping["Udp"] {
		mapping = append(mapping, fmt.Sprintf("%s->%s/udp", public, private))
	}
	sort.Strings(mapping)
	return strings.Join(mapping, ", ")
}

type PortMapping map[string]string
type HostConfig struct {
	Binds           []string
	ContainerIDFile string
}

