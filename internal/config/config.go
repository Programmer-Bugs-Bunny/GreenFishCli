package config

// ProjectConfig 项目配置结构
type ProjectConfig struct {
	ProjectName string // 项目名称
	ModuleName  string // Go 模块名
	Description string // 项目描述
	Author      string // 作者
	Username    string // GitHub 用户名
	RepoName    string // 仓库名
	Port        string // 服务端口
	Timezone    string // 时区
}

// GetTemplateVars 获取模板变量映射
func (c *ProjectConfig) GetTemplateVars() map[string]string {
	return map[string]string{
		"{{.ProjectName}}": c.ProjectName,
		"{{.ModuleName}}":  c.ModuleName,
		"{{.Description}}": c.Description,
		"{{.Author}}":      c.Author,
		"{{.Username}}":    c.Username,
		"{{.RepoName}}":    c.RepoName,
		"{{.Port}}":        c.Port,
		"{{.Timezone}}":    c.Timezone,
	}
}
