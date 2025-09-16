package template

import (
	"fmt"
	"text/template"
)

type Manager struct {
	templates map[string]*template.Template
}

func New() *Manager {
	return &Manager{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplate 加载模板
func (m *Manager) LoadTemplate(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return fmt.Errorf("解析模板 %s 失败: %w", name, err)
	}

	m.templates[name] = tmpl
	return nil
}

// GetTemplate 获取模板
func (m *Manager) GetTemplate(name string) (*template.Template, bool) {
	tmpl, exists := m.templates[name]
	return tmpl, exists
}
