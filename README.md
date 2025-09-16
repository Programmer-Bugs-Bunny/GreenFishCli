# GoWeb CLI

🚀 基于 [GreenFishCli](https://github.com/Programmer-Bugs-Bunny/GreenFishCli) 的 Go Web 项目脚手架工具

类似于 Vue CLI，可以快速创建标准化的 Go Web 项目，内置 Gin、GORM、Zap 日志等常用组件。

## ✨ 特性

- 🎯 **交互式创建**: 类似 `vue create`，提供友好的命令行交互
- 🛠️ **自动配置**: 自动替换模块名、配置文件和项目信息
- 📦 **完整模板**: 基于成熟的 GoWebTemplate，包含最佳实践
- 🔧 **高度定制**: 支持自定义项目名、作者、端口、时区等
- 🚀 **即开即用**: 生成的项目可直接运行

## 📦 安装

### 方法一：一键安装（推荐）

```bash
curl -sSL https://raw.githubusercontent.com/Programmer-Bugs-Bunny/GreenFishCli/main/install.sh | bash
```

### 方法二：手动编译安装

```bash
# 克隆项目
git clone https://github.com/Programmer-Bugs-Bunny/goweb-cli.git
cd goweb-cli

# 编译安装
go build -o goweb .
sudo mv goweb /usr/local/bin/
```

### 方法三：Go install（开发中）

```bash
go install github.com/Programmer-Bugs-Bunny/goweb-cli@latest
```

## 🚀 使用方法

### 创建新项目

```bash
# 交互式创建项目
goweb create

# 指定项目名创建
goweb create my-awesome-project

# 查看帮助
goweb --help
goweb create --help
```

### 交互式配置

运行 `goweb create` 后，工具会引导您配置：

- **项目名称**: 作为目录名和默认模块名
- **Go 模块名**: 完整的模块路径，如 `github.com/username/project`
- **项目描述**: 用于 README 和文档
- **作者名称**: 用于许可证和文档
- **服务端口**: 默认 8080
- **时区设置**: 支持多个常用时区

### 示例输出

```bash
$ goweb create my-web-app

🚀 GoWeb CLI - 项目脚手架工具
===============================

? 项目名称: my-web-app
? Go 模块名: github.com/myname/my-web-app
? 项目描述: 基于 GoWebTemplate 创建的 my-web-app 项目
? 作者: Your Name
? 服务端口: 8080
? 时区: Asia/Shanghai

🚀 正在创建项目 'my-web-app'...
📁 模板源目录: /path/to/template
📁 目标目录: my-web-app
✅ 模板复制完成
🔄 处理模板变量替换...
✅ 已处理: go.mod
✅ 已处理: README.md
✅ 已处理: config/app.yaml
✅ 已处理: atlas.hcl
✅ 项目 'my-web-app' 创建成功！

📋 后续步骤:
   cd my-web-app
   go mod tidy
   # 修改 config/app.yaml 中的数据库配置
   go run main.go
```

## 📁 生成的项目结构

```
my-web-app/
├── config/         # 配置文件
├── database/       # 数据库初始化
├── middlewares/    # 中间件
├── models/         # 数据模型
├── modules/        # 业务模块
├── routes/         # 路由注册
├── utils/          # 工具方法
├── cmd/            # 命令行工具
├── atlas.hcl       # 数据库迁移配置
├── main.go         # 程序入口
├── go.mod          # 已配置正确的模块名
└── README.md       # 项目说明
```

## 🔧 自动处理的内容

CLI 工具会自动处理以下内容：

1. **go.mod**: 替换为指定的模块名
2. **README.md**: 更新项目信息和说明
3. **config/app.yaml**: 配置端口和时区
4. **atlas.hcl**: 数据库迁移配置
5. **导入路径**: 所有 Go 文件中的导入路径自动更新

## 🎛️ 命令行选项

```bash
# 创建项目的各种方式
goweb create                                    # 交互式创建
goweb create my-project                         # 指定项目名
goweb create my-project --module github.com/me/my-project  # 指定模块名
goweb create my-project --author "Your Name"   # 指定作者
goweb create my-project --port 3000            # 指定端口

# 其他命令
goweb --version                                 # 查看版本
goweb --help                                    # 查看帮助
```

## 🛠️ 开发

```bash
# 克隆项目
git clone https://github.com/Programmer-Bugs-Bunny/goweb-cli.git
cd goweb-cli

# 安装依赖
go mod tidy

# 运行
go run . create test-project

# 构建
go build -o goweb .
```

## 📋 TODO

- [ ] 支持多种项目模板（基础版、完整版、微服务版）
- [ ] 添加代码生成器功能（controller、service、model）
- [ ] 支持插件系统
- [ ] 添加项目升级功能
- [ ] 支持配置文件预设

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
