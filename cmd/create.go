package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/config"
	"github.com/Programmer-Bugs-Bunny/GreenFishCli/internal/generator"
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "创建新的 Go Web 项目",
	Long: `基于 GoWebTemplate 模板创建新的 Go Web 项目。
支持交互式配置项目参数，自动替换模块名和配置文件。`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string

		if len(args) > 0 {
			projectName = args[0]
		} else {
			prompt := &survey.Input{
				Message: "项目名称:",
				Help:    "将作为目录名和 Go 模块名的一部分",
			}
			survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required))
		}

		if projectName == "" {
			fmt.Println("❌ 项目名称不能为空")
			os.Exit(1)
		}

		// 检查目录是否已存在
		if _, err := os.Stat(projectName); !os.IsNotExist(err) {
			fmt.Printf("❌ 目录 '%s' 已存在\n", projectName)
			os.Exit(1)
		}

		// 交互式收集项目配置
		projectConfig, err := collectProjectConfig(projectName)
		if err != nil {
			fmt.Printf("❌ 配置收集失败: %v\n", err)
			os.Exit(1)
		}

		// 生成项目
		gen := generator.New()
		if err := gen.Generate(projectConfig); err != nil {
			fmt.Printf("❌ 项目生成失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ 项目 '%s' 创建成功！\n", projectName)
		fmt.Println("\n📋 后续步骤:")
		fmt.Printf("   cd %s\n", projectName)
		fmt.Println("   go mod tidy")
		fmt.Println("   # 修改 config/app.yaml 中的数据库配置")
		fmt.Println("   go run main.go")
	},
}

func collectProjectConfig(projectName string) (*config.ProjectConfig, error) {
	cfg := &config.ProjectConfig{
		ProjectName: projectName,
	}

	questions := []*survey.Question{
		{
			Name: "module",
			Prompt: &survey.Input{
				Message: "Go 模块名:",
				Default: fmt.Sprintf("github.com/yourusername/%s", projectName),
				Help:    "完整的 Go 模块路径，如: github.com/username/project",
			},
			Validate: survey.Required,
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "项目描述:",
				Default: fmt.Sprintf("基于 GoWebTemplate 创建的 %s 项目", projectName),
			},
		},
		{
			Name: "author",
			Prompt: &survey.Input{
				Message: "作者:",
				Default: "Your Name",
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "服务端口:",
				Default: "8080",
			},
		},
		{
			Name: "timezone",
			Prompt: &survey.Select{
				Message: "时区:",
				Options: []string{
					"Asia/Shanghai",
					"Asia/Ho_Chi_Minh",
					"UTC",
					"America/New_York",
					"Europe/London",
				},
				Default: "Asia/Shanghai",
			},
		},
	}

	answers := struct {
		Module      string `survey:"module"`
		Description string `survey:"description"`
		Author      string `survey:"author"`
		Port        string `survey:"port"`
		Timezone    string `survey:"timezone"`
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return nil, err
	}

	cfg.ModuleName = answers.Module
	cfg.Description = answers.Description
	cfg.Author = answers.Author
	cfg.Port = answers.Port
	cfg.Timezone = answers.Timezone

	// 从模块名提取用户名和仓库名
	parts := strings.Split(answers.Module, "/")
	if len(parts) >= 2 {
		cfg.Username = parts[len(parts)-2]
		cfg.RepoName = parts[len(parts)-1]
	}

	return cfg, nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("module", "m", "", "Go 模块名")
	createCmd.Flags().StringP("author", "a", "", "作者名称")
	createCmd.Flags().StringP("port", "p", "8080", "服务端口")
}
