package glslgen

type Variable struct {
	Type      string
	Precision string
	Name      string
}

type Makro struct {
	Name  string
	Value string
}

type Function struct {
	Prototype string
	Body      string
}

type Module struct {
	Uniforms  []Variable
	Functions []Function
	Name      string
	Body      string
}

type Generator struct {
	Version string
	Outputs []Variable
	Makros  []Makro
	Globals []Variable
	Modules []Module
}

type VertexGenerator struct {
	Generator
	Attributes []Variable
}

type FragmentGenerator struct {
	Generator
	Inputs []Variable
}
