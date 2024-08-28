package cmadapter

type Adapter interface {
	GetValue(string) (string, error)
	MonitorToReLoad(msg []byte)
}
