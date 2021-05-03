package dapmesh

type DAPMeshCall struct {
	Message    string
	Subscriber []string
	Groups     []string
	Emergency  bool
}

type DAPMeshMessage struct {
	Version int
	Type    string
	Payload interface{}
}
