package service

type Service struct {
	Name    string
	Image   string
	Port    int
	Enabled bool
	Env     map[string]string
}

type Status struct {
	Name    string
	Running bool
	State   string
	Port    int
	Image   string
}

type Repository interface {
	GetAll() ([]Service, error)
	GetStatus(name string) (Status, error)
	Start(name string) error
	Stop(name string) error
	Restart(name string) error
	StartAll() error
	StopAll() error
	RestartAll() error
}
