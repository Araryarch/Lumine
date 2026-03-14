package project

type Type string

const (
	TypePHP       Type = "PHP"
	TypeLaravel   Type = "Laravel"
	TypeWordPress Type = "WordPress"
	TypeNodeJS    Type = "Node.js"
	TypeStatic    Type = "Static"
	TypeUnknown   Type = "Unknown"
)

type Project struct {
	Name string
	Type Type
	Path string
	URL  string
}

type Repository interface {
	List() ([]Project, error)
	Create(name string, projectType Type, phpVersion string) error
	Delete(name string) error
	GetType(path string) Type
}

type Service interface {
	CreatePHP(name, phpVersion string) error
	CreateLaravel(name, phpVersion string) error
	CreateWordPress(name, phpVersion string) error
	CreateNodeJS(name string) error
	CreateStatic(name string) error
}
