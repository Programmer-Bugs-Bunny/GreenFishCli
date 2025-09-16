package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/config"
	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/template"
)

type Generator struct {
	templateManager *template.Manager
}

func New() *Generator {
	return &Generator{
		templateManager: template.New(),
	}
}

// Generate 生成项目
func (g *Generator) Generate(cfg *config.ProjectConfig) error {
	fmt.Printf("🚀 正在创建项目 '%s'...\n", cfg.ProjectName)

	// 创建项目目录
	if err := os.MkdirAll(cfg.ProjectName, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w", err)
	}

	// 获取模板源目录（当前项目的父目录，排除cli-tool目录）
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}

	// 如果当前在 cli-tool 目录中，则模板源为上级目录
	templateSourceDir := filepath.Dir(currentDir)
	if !strings.HasSuffix(currentDir, "cli-tool") {
		templateSourceDir = currentDir
	}

	fmt.Printf("📁 模板源目录: %s\n", templateSourceDir)
	fmt.Printf("📁 目标目录: %s\n", cfg.ProjectName)

	// 复制模板文件
	if err := g.copyTemplate(templateSourceDir, cfg.ProjectName, cfg); err != nil {
		return fmt.Errorf("复制模板失败: %w", err)
	}

	// 处理模板变量替换
	if err := g.processTemplates(cfg); err != nil {
		return fmt.Errorf("处理模板变量失败: %w", err)
	}

	fmt.Println("✅ 模板复制完成")
	return nil
}

// copyTemplate 复制模板文件
func (g *Generator) copyTemplate(sourceDir, targetDir string, cfg *config.ProjectConfig) error {
	// 需要跳过的目录和文件
	skipDirs := map[string]bool{
		"cli-tool":   true,
		".git":       true,
		".idea":      true,
		"logs":       true,
		"migrations": true,
	}

	skipFiles := map[string]bool{
		".DS_Store": true,
	}

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 检查是否需要跳过
		pathParts := strings.Split(relPath, string(filepath.Separator))
		for _, part := range pathParts {
			if skipDirs[part] || skipFiles[part] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return g.copyFile(path, targetPath)
	})
}

// copyFile 复制单个文件
func (g *Generator) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// processTemplates 处理模板变量替换
func (g *Generator) processTemplates(cfg *config.ProjectConfig) error {
	fmt.Println("🔄 处理模板变量替换...")

	// 需要进行变量替换的文件
	templateFiles := []string{
		"go.mod",
		"README.md",
		"config/app.yaml",
		"atlas.hcl",
	}

	templateVars := cfg.GetTemplateVars()

	for _, file := range templateFiles {
		filePath := filepath.Join(cfg.ProjectName, file)
		if err := g.replaceInFile(filePath, templateVars); err != nil {
			// 文件可能不存在，继续处理其他文件
			fmt.Printf("⚠️  处理文件 %s 时出错: %v\n", file, err)
			continue
		}
		fmt.Printf("✅ 已处理: %s\n", file)
	}

	return nil
}

// replaceInFile 在文件中进行变量替换
func (g *Generator) replaceInFile(filePath string, vars map[string]string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 进行变量替换
	newContent := string(content)
	for placeholder, value := range vars {
		newContent = strings.ReplaceAll(newContent, placeholder, value)
	}

	// 特殊处理：替换模板仓库的模块名为新的模块名
	newContent = strings.ReplaceAll(newContent, "github.com/Programmer-Bugs-Bunny/GoWebTemplate", vars["{{.ModuleName}}"])

	// 写回文件
	return os.WriteFile(filePath, []byte(newContent), 0644)
}
