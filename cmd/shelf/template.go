package main

import (
	"os"
	"text/template"
)

type Import struct {
	Name string
	Path string
}

type LoadMethods struct {
	Name string
	Flag int
}

type Repository struct {
	Name            string
	InterfaceName   string
	TypeName        string
	ConstructorName string
	LoadMethods     []LoadMethods
}

type Template struct {
	PackageName  string
	Imports      []Import
	Repositories []Repository
}

const (
	shelfFileTemplate = `{package {{ .PackageName}}

import ({{ range $import := .Imports }} {{if $import.Name }}
	{{$import.Name}} "{{$import.Path}}" {{else}} 
	"{{$import.Path}}" {{end}} {{end}}
)

func init() { {{range $repository := .Repositories }}
	shelf.RegisterRepository("{{$repository.Name}}", shelf__new{{$repository.ConstructorName}}()){{end}}
}

{{range $repository := .Repositories -}}
type shelf__{{$repository.TypeName}} struct {
	ParentRepository {{$repository.InterfaceName}}
	LoadFlags	int
}

func shelf__new{{$repository.ConstructorName}}() *shelf__{{$repository.TypeName}} {
	return &shelf__{{$repository.TypeName}}{}
}

{{range $loadMethod := $repository.LoadMethods -}}
func (r *shelf__{{$repository.TypeName}}) {{$loadMethod.Name}}() {{$repository.InterfaceName}} {
	if r.ParentRepository == nil {
		clone := &shelf__{{$repository.TypeName}}{
			ParentRepository : r,
		}
		clone.LoadFlags |= {{$loadMethod.Flag}}
		return clone
	}

	r.LoadFlags |= {{$loadMethod.Flag}}
	return r
}

{{end}}
{{end}}
`
)

func GenerateCode() {
	x, _ := template.New("shelfFileTemplate").Parse(shelfFileTemplate)
	x.Execute(os.Stdout, Template{
		PackageName: "auto_generated",
		Imports: []Import{
			{Name: "", Path: "context"},
			{Name: "", Path: "github.com/procyon-projects/shelf"},
			{Name: "", Path: "github.com/test"}},
		Repositories: []Repository{
			{
				Name:            "user-repository",
				InterfaceName:   "UserRepository",
				TypeName:        "userRepositoryImpl",
				ConstructorName: "UserRepositoryImpl",
				LoadMethods: []LoadMethods{
					{Name: "LoadPost", Flag: 1},
					{Name: "LoadCreditCart", Flag: 2},
				},
			},
			{
				Name:            "post-repository",
				InterfaceName:   "PostRepository",
				TypeName:        "postRepositoryImpl",
				ConstructorName: "PostRepositoryImpl",
			},
		},
	})
}
