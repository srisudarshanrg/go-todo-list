package models

type TemplateData struct {
	Errors  map[string]string
	Info    map[string]string
	Success map[string]string
	Data    map[string]interface{}
}
