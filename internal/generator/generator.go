package generator

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/config"
	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/template"
)

const (
	// 模板仓库地址
	templateRepoURL = "https://github.com/Programmer-Bugs-Bunny/GreenFish.git"
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

	// 克隆模板仓库（git clone会自动创建目标目录）
	fmt.Printf("📁 目标目录: %s\n", cfg.ProjectName)
	fmt.Printf("📦 正在克隆模板仓库...\n")
	
	if err := g.cloneTemplate(cfg.ProjectName); err != nil {
		return fmt.Errorf("克隆模板失败: %w", err)
	}

	// 处理模板变量替换
	if err := g.processTemplates(cfg); err != nil {
		return fmt.Errorf("处理模板变量失败: %w", err)
	}

	// 处理 Go 文件中的导入路径替换
	if err := g.processGoFiles(cfg); err != nil {
		return fmt.Errorf("处理 Go 文件导入路径失败: %w", err)
	}

	fmt.Println("✅ 模板复制完成")
	return nil
}

// cloneTemplate 克隆模板仓库
func (g *Generator) cloneTemplate(targetDir string) error {
	// 使用 git clone 克隆模板仓库到目标目录
	cmd := exec.Command("git", "clone", templateRepoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 git clone 失败: %w", err)
	}
	
	// 删除克隆后的 .git 目录，使其成为一个干净的项目
	gitDir := filepath.Join(targetDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		// 删除 .git 目录失败不是致命错误，继续执行
		fmt.Printf("⚠️  删除 .git 目录失败: %v\n", err)
	}
	
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
	newContent = strings.ReplaceAll(newContent, "github.com/Programmer-Bugs-Bunny/GreenFish", vars["{{.ModuleName}}"])
	
	// 替换所有的 go-web-template 导入路径为新的模块名
	newContent = strings.ReplaceAll(newContent, "go-web-template", vars["{{.ModuleName}}"])

	// 写回文件
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// processGoFiles 处理 Go 文件中的导入路径替换
func (g *Generator) processGoFiles(cfg *config.ProjectConfig) error {
	fmt.Println("🔄 处理 Go 文件导入路径...")
	
	projectPath := cfg.ProjectName
	goFileCount := 0
	
	// 遍历项目目录，找到所有 .go 文件
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 跳过不需要处理的目录
		if info.IsDir() {
			dirName := filepath.Base(path)
			if dirName == ".git" || dirName == "GreenFishCli" || dirName == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		
		// 只处理 .go 文件
		if filepath.Ext(path) != ".go" {
			return nil
		}
		
		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", path, err)
		}
		
		// 替换导入路径
		newContent := string(content)
		oldImportPath := "go-web-template"
		newImportPath := cfg.ModuleName
		
		if strings.Contains(newContent, oldImportPath) {
			newContent = strings.ReplaceAll(newContent, oldImportPath, newImportPath)
			
			// 写回文件
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("写入文件 %s 失败: %w", path, err)
			}
			
			// 计算相对路径用于显示
			relPath, _ := filepath.Rel(projectPath, path)
			fmt.Printf("✅ 已处理: %s\n", relPath)
			goFileCount++
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	fmt.Printf("🔄 共处理了 %d 个 Go 文件\n", goFileCount)
	return nil
}
