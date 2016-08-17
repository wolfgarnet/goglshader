package goglshader

import (
	"fmt"
	"github.com/go-gl/gl/v4.5-core/gl"
	"io/ioutil"
	"os"
	"strings"
)

func LoadGLProgram(vertexShaderFile, fragmentShaderFile string) (uint32, error) {
	vertexShader, err := LoadShader(vertexShaderFile, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := LoadShader(fragmentShaderFile, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	return CreateProgram(vertexShader, fragmentShader)
}

func InitializeProgram(vertexShaderString, fragmentShaderString string) (uint32, error) {
	vertexShaderString += string('\x00')
	vertexShader, err := CompileShader(vertexShaderString, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShaderString += string('\x00')
	fragmentShader, err2 := CompileShader(fragmentShaderString, gl.FRAGMENT_SHADER)
	if err2 != nil {
		return 0, err2
	}

	return CreateProgram(vertexShader, fragmentShader)
}

func CreateProgram(vertexShader, fragmentShader uint32) (uint32, error) {
	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func LoadShader(shaderFile string, shaderType uint32) (uint32, error) {
	file, err := os.Open(shaderFile)
	if err != nil {
		return 0, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	content = append(content, '\x00')
	return CompileShader(string(content), shaderType)
}


func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}


	return shader, nil
}

