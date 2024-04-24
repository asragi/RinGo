package router

import (
	"fmt"
	"regexp"
	"strings"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

func NewMethod(method string) (Method, error) {
	switch method {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "DELETE":
		return DELETE, nil
	case "PATCH":
		return PATCH, nil
	default:
		return "", fmt.Errorf("invalid method: %s", method)
	}
}

func (m Method) Is(method Method) bool {
	return string(m) == string(method)
}

func (m Method) ToString() string {
	return string(m)
}

type Path struct {
	path string
}

func (p *Path) ToString() string {
	return p.path
}

func NewPathData(path string) (*Path, error) {
	isValid := func(path string) bool {
		return true
	}(path)
	if !isValid {
		return nil, fmt.Errorf("invalid path: %s", path)
	}
	return &Path{path: path}, nil
}

func (p *Path) GetStringByIndex(index int) (string, error) {
	pathArray := p.toPathArray()
	if index < 0 || index >= len(pathArray) {
		return "", fmt.Errorf("index out of range: %d", index)
	}
	return pathArray[index], nil
}

type paramExpression string

func (p paramExpression) Match(s string) bool {
	return string(p) == s
}

type SamplePath string

func NewSamplePath(path string) (SamplePath, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("invalid path: %s", path)
	}
	return SamplePath(path), nil
}

func (s SamplePath) GetParameterIndex(paramExpression paramExpression) (int, error) {
	pathArray := s.toPathArray()
	for i, v := range pathArray {
		if paramExpression.Match(v) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("paramExpression not found: %s", paramExpression)
}

type PathMatchPattern struct {
	exp *regexp.Regexp
}

func NewPathMatchPattern(samplePath SamplePath) (*PathMatchPattern, error) {
	const idPattern string = `[a-zA-Z0-9\-]+`
	pathSplit := samplePath.toPathArray()
	patterns := make([]string, len(pathSplit)-1)
	for i := range pathSplit {
		if i == 0 {
			continue
		}
		targetString := pathSplit[i]
		if len(targetString) == 0 {
			return nil, fmt.Errorf("invalid path: %s", samplePath)
		}
		if firstCharacter := targetString[0]; firstCharacter == '{' {
			patterns[i-1] = idPattern
			continue
		}
		patterns[i-1] = targetString
	}
	regexpPattern := func() string {
		builder := strings.Builder{}
		builder.WriteString("^")
		for _, v := range patterns {
			builder.WriteString("/")
			builder.WriteString(v)
		}
		builder.WriteString("$")
		return builder.String()
	}()
	compiledRegExp, err := regexp.Compile(regexpPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", samplePath)
	}
	return &PathMatchPattern{exp: compiledRegExp}, nil
}

func (p *PathMatchPattern) Match(path *Path) bool {
	return p.exp.MatchString(path.path)
}

func (s SamplePath) toPathArray() []string {
	return toPathArray(string(s))
}

func (p *Path) toPathArray() []string {
	return toPathArray(p.path)
}

func toPathArray(path string) []string {
	pathSplit := strings.Split(path, "/")
	return pathSplit
}
