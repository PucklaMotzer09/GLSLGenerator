package glslgen

import (
	"strconv"
	"strings"
)

func (this *Variable) String(version uint16) string {
	var p = ""
	if isES {
		p = this.Precision + " "
	}
	tiep := this.Type
	if tiep == "sampler2DMS" {
		if isES && version < 320 {
			tiep = "sampler2D"
		}
	}
	return p + tiep + " " + this.Name + ";"
}

func (this *Variable) UniformString(version uint16) string {
	return "uniform " + this.String(version)
}

func (this *Variable) AttributeString(index, version uint16) string {
	if version < 150 || (isES && version < 300) {
		return "attribute " + this.String(version)
	} else {
		return "in " + this.String(version)
	}
}

func (this *Variable) OutputString(version uint16) string {
	if version < 150 || (isES && version < 300) {
		return "varying " + this.String(version)
	} else {
		return "out " + this.String(version)
	}
}

func (this *Variable) InputString(version uint16) string {
	if version < 150 || (isES && version < 300) {
		return "varying " + this.String(version)
	} else {
		return "in " + this.String(version)
	}
}

func (this *Makro) String() string {
	return "#define " + this.Name + " " + this.Value
}

func (this *Struct) String(version uint16) (str string) {
	str += "struct " + this.Name + "\n{"
	for _, v := range this.Variables {
		str += v.String(version) + "\n"
	}
	str += "};"
	return
}

func (this *Function) String() string {
	return this.Prototype + "\n{\n" + this.Body + "\n}\n"
}

func (this *Function) PrototypeString() string {
	return this.Prototype + ";"
}

func (this *Module) UniformsString(version uint16) (str string) {
	for _, u := range this.Uniforms {
		str += u.UniformString(version) + "\n"
	}
	return
}

func (this *Module) FunctionPrototypesString() (str string) {
	for _, f := range this.Functions {
		str += f.PrototypeString() + "\n"
	}
	return
}

func (this *Module) FunstionsString() (str string) {
	for _, f := range this.Functions {
		str += f.String() + "\n"
	}
	return
}

func (this *Module) PrototypeString(index uint8) string {
	if this.Body == "" {
		return ""
	}
	if this.Name != "" {
		return "void " + this.Name + "();"
	} else {
		return "void module" + strconv.FormatUint(uint64(index), 10) + "();"
	}
}

func (this *Module) CallString(index uint8) string {
	if this.Body == "" {
		return ""
	}
	var name string
	if this.Name != "" {
		name = this.Name
	} else {
		name = "module" + strconv.FormatUint(uint64(index), 10)
	}
	return name + "();"
}

func (this *Module) String(index uint8) string {
	if this.Body == "" {
		return ""
	}
	ps := this.PrototypeString(index)
	ps = ps[:len(ps)-1]

	return ps + "\n{\n" + this.Body + "\n}"
}

func (this *Generator) String() string {
	return ""
}

func printPrecisions() (str string) {
	str += "\n"
	str += "precision highp float;\n"
	str += "precision highp sampler2D;\n"
	str += "\n"
	return
}

func (this *VertexGenerator) String() (str string) {
	var temp uint64
	var versioni uint16
	if this.Version == "WebGL" {
		isES = true
		isWebGL = true
	} else {
		if this.Version == "100" || (len(this.Version) == 6 && this.Version[4:] == "es") {
			temp, _ = strconv.ParseUint(this.Version[:3], 10, 16)
			isES = true
		} else {
			temp, _ = strconv.ParseUint(this.Version, 10, 16)
		}
		versioni = uint16(temp)
	}
	versioni = uint16(temp)

	if !isWebGL {
		str += "#version " + this.Version + "\n\n"
	}

	for _, m := range this.Makros {
		str += m.String() + "\n"
	}
	if len(this.Makros) != 0 {
		str += "\n"
	}

	if isES {
		str += printPrecisions()
	}

	for i, a := range this.Attributes {
		str += a.AttributeString(uint16(i), versioni) + "\n"
	}
	str += "\n"

	for _, o := range this.Outputs {
		str += o.OutputString(versioni) + "\n"
	}
	if len(this.Outputs) != 0 {
		str += "\n"
	}

	var hasStructs = false
	for _, m := range this.Modules {
		for _, s := range m.Structs {
			hasStructs = true
			str += s.String(versioni) + "\n"
		}
	}
	if hasStructs {
		str += "\n"
	}

	var hasUniforms = false
	for _, m := range this.Modules {
		for _, u := range m.Uniforms {
			hasUniforms = true
			str += u.UniformString(versioni) + "\n"
		}
	}
	if hasUniforms {
		str += "\n"
	}

	for _, g := range this.Globals {
		str += g.String(versioni)
	}
	if len(this.Globals) != 0 {
		str += "\n"
	}

	var hasFunctions = false
	for _, m := range this.Modules {
		for _, f := range m.Functions {
			hasFunctions = true
			str += f.PrototypeString() + "\n"
		}
	}
	if hasFunctions {
		str += "\n"
	}

	for i, m := range this.Modules {
		str += m.PrototypeString(uint8(i)) + "\n"
	}
	str += "\n"

	str += "void main()\n{\n"
	for i, m := range this.Modules {
		str += m.CallString(uint8(i)) + "\n"
	}
	str += "}\n\n"

	for i, m := range this.Modules {
		str += m.String(uint8(i)) + "\n\n"
	}
	str += "\n"

	for _, m := range this.Modules {
		for _, f := range m.Functions {
			str += f.String() + "\n\n"
		}
	}
	if hasFunctions {
		str += "\n"
	}

	return
}

func (this *FragmentGenerator) String() (str string) {

	var temp uint64
	var versioni uint16
	if this.Version == "WebGL" {
		isES = true
		isWebGL = true
	} else {
		if this.Version == "100" || (len(this.Version) == 6 && this.Version[4:] == "es") {
			temp, _ = strconv.ParseUint(this.Version[:3], 10, 16)
			isES = true
		} else {
			temp, _ = strconv.ParseUint(this.Version, 10, 16)
		}
		versioni = uint16(temp)
	}

	if isES && versioni >= 300 {
		this.AddOutput(Variable{"vec4", "highp", "glFragColor"})
	}

	if !isWebGL {
		str += "#version " + this.Version + "\n\n"
	}

	for _, m := range this.Makros {
		str += m.String() + "\n"
	}
	if len(this.Makros) != 0 {
		str += "\n"
	}

	if isES {
		str += printPrecisions()
	}

	for _, i := range this.Inputs {
		str += i.InputString(versioni) + "\n"
	}
	if len(this.Inputs) != 0 {
		str += "\n"
	}

	for _, o := range this.Outputs {
		str += o.OutputString(versioni) + "\n"
	}
	if len(this.Outputs) != 0 {
		str += "\n"
	}

	var hasStructs = false
	for _, m := range this.Modules {
		for _, s := range m.Structs {
			hasStructs = true
			str += s.String(versioni) + "\n"
		}
	}
	if hasStructs {
		str += "\n"
	}

	var hasUniforms = false
	for _, m := range this.Modules {
		for _, u := range m.Uniforms {
			hasUniforms = true
			str += u.UniformString(versioni) + "\n"
		}
	}
	if hasUniforms {
		str += "\n"
	}

	for _, g := range this.Globals {
		str += g.String(versioni) + "\n"
	}
	if len(this.Globals) != 0 {
		str += "\n"
	}

	var hasFunctions = false
	for _, m := range this.Modules {
		for _, f := range m.Functions {
			hasFunctions = true
			str += f.PrototypeString() + "\n"
		}
	}
	if hasFunctions {
		str += "\n"
	}

	for i, m := range this.Modules {
		str += m.PrototypeString(uint8(i)) + "\n"
	}
	str += "\n"

	str += "void main()\n{\n"
	for i, m := range this.Modules {
		str += m.CallString(uint8(i)) + "\n"
	}
	str += "}\n\n"

	for i, m := range this.Modules {
		str += m.String(uint8(i)) + "\n\n"
	}
	str += "\n"

	for _, m := range this.Modules {
		for _, f := range m.Functions {
			str += f.String() + "\n\n"
		}
	}
	if hasFunctions {
		str += "\n"
	}

	if isES && versioni >= 300 {
		str = strings.Replace(str, "gl_FragColor", "glFragColor", -1)
	}

	return
}
