package models

type OpenAPI3 struct {
	OpenAPI string              `yaml:"openapi"`
	Info    Info                `yaml:"info"`
	Servers []Server            `yaml:"servers"`
	Paths   map[string]PathItem `yaml:"paths"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Server struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

type PathItem struct {
	Get     *Operation `yaml:"get,omitempty"`
	Post    *Operation `yaml:"post,omitempty"`
	Put     *Operation `yaml:"put,omitempty"`
	Delete  *Operation `yaml:"delete,omitempty"`
	Patch   *Operation `yaml:"patch,omitempty"`
	Options *Operation `yaml:"options,omitempty"`
	Head    *Operation `yaml:"head,omitempty"`
}

type Operation struct {
	Summary     string              `yaml:"summary,omitempty"`
	Description string              `yaml:"description,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty"`
	RequestBody *RequestBody        `yaml:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses"`
	Security    []Security          `yaml:"security,omitempty"`
}

type Parameter struct {
	Name        string  `yaml:"name"`
	In          string  `yaml:"in"`
	Description string  `yaml:"description,omitempty"`
	Required    bool    `yaml:"required,omitempty"`
	Schema      *Schema `yaml:"schema,omitempty"`
}

type RequestBody struct {
	Description string             `yaml:"description,omitempty"`
	Content     map[string]Content `yaml:"content"`
}

type Response struct {
	Description string             `yaml:"description"`
	Content     map[string]Content `yaml:"content,omitempty"`
}

type Content struct {
	Schema *Schema `yaml:"schema"`
}

type Schema struct {
	Type       string            `yaml:"type,omitempty"`
	Items      *Items            `yaml:"items,omitempty"`
	Properties map[string]Schema `yaml:"properties,omitempty"`
}

type Items struct {
	Type string `yaml:"type,omitempty"`
}

type Security struct {
	ApiKeyAuth []string `yaml:"ApiKeyAuth,omitempty"`
}
