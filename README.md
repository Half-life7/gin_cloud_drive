# gin_cloud_drive

一个基于Go语言和Gin框架开发的个人网盘系统，提供文件上传、下载、管理、系统监控等功能。

[![GitHub stars](https://img.shields.io/github/stars/Half-life7/gin_cloud_drive?style=social)](https://github.com/Half-life7/gin_cloud_drive/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/Half-life7/gin_cloud_drive?style=social)](https://github.com/Half-life7/gin_cloud_drive/network/members)
[![GitHub license](https://img.shields.io/github/license/Half-life7/gin_cloud_drive)](https://github.com/Half-life7/gin_cloud_drive/blob/main/LICENSE)

## 项目地址

GitHub: [https://github.com/Half-life7/gin_cloud_drive](https://github.com/Half-life7/gin_cloud_drive)

## 功能特性

### 文件管理
- ✅ 文件上传（支持拖拽上传）
- ✅ 文件下载
- ✅ 文件预览
- ✅ 文件移动
- ✅ 文件删除
- ✅ 新建文件夹
- ✅ 手动上传（文件缓冲区域，选择路径后再上传）
- ✅ 支持批量上传

### 系统功能
- ✅ 系统状态监控（CPU、内存、磁盘、网络）
- ✅ 系统状态趋势图
- ✅ 日志管理（操作日志、访问日志）
- ✅ 暗色/亮色模式切换
- ✅ 轮播图功能

### 用户体验
- ✅ 响应式设计
- ✅ 友好的错误提示
- ✅ 上传进度显示
- ✅ 主题持久化存储

## 技术栈

### 后端
- **语言**：Go 1.18+
- **框架**：Gin Web Framework
- **系统监控**：gopsutil
- **日志管理**：自定义日志系统

### 前端
- **HTML5**
- **CSS3**（使用CSS变量实现主题切换）
- **JavaScript**（原生JS）
- **图表库**：Chart.js

### 数据存储
- **文件存储**：本地文件系统
- **日志存储**：本地文件

## 安装和运行

### 环境要求
- Go 1.18+
- 支持的操作系统：Windows、Linux、macOS

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/Half-life7/gin_cloud_drive.git
   cd gin_cloud_drive
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **运行项目**
   ```bash
   go run main.go
   ```

4. **访问系统**
   打开浏览器，访问 `http://localhost:8080`

5. **构建项目**
   ```bash
   go build -o gin_cloud_drive.exe
   ```
   或在Linux/macOS上：
   ```bash
   go build -o gin_cloud_drive
   ```

## 使用说明

### 文件上传
1. **选择文件**：点击"选择文件"按钮或拖拽文件到上传区域
2. **查看待上传文件**：文件会显示在缓冲区域，可以查看文件名和大小
3. **选择上传路径**：点击"选择文件夹"按钮，选择文件要上传到的路径
4. **开始上传**：点击"开始上传"按钮，文件会被上传到指定路径

### 文件管理
- **进入文件夹**：点击文件夹的"进入"按钮
- **下载文件**：点击文件的"下载"按钮
- **预览文件**：点击文件的"预览"按钮
- **移动文件**：点击文件的"移动"按钮，选择目标路径
- **删除文件**：点击文件的"删除"按钮

### 系统状态
- 点击导航栏的"关于" -> "系统状态"，查看系统CPU、内存、磁盘、网络等信息
- 系统状态会实时更新，并显示趋势图

### 日志管理
- 点击导航栏的"关于" -> "日志管理"，查看系统操作日志和访问日志
- 可以根据时间、操作类型、IP地址等条件查询日志

### 主题切换
- 点击导航栏右上角的太阳/月亮图标，切换亮色/暗色模式
- 主题设置会保存在本地存储中，下次访问时会自动应用

## 项目结构

```
gin_cloud_drive/
├── backend/
│   ├── controllers/      # 控制器
│   ├── middleware/       # 中间件
│   └── routes/           # 路由
├── frontend/             # 前端代码
│   ├── css/              # CSS样式
│   ├── js/               # JavaScript代码
│   ├── about.html        # 关于页面
│   ├── disk.html         # 网盘页面
│   ├── index.html        # 首页
│   ├── logs.html         # 日志页面
│   └── system.html       # 系统状态页面
├── logger/               # 日志系统
├── system/               # 系统监控
├── upload/               # 文件上传目录
├── main.go               # 主程序
├── go.mod                # Go模块依赖
├── go.sum                # Go模块校验和
└── README.md             # 项目文档
```

## 配置说明

### 端口配置
- 默认端口：8080
- 可以在`main.go`文件中修改端口号

### 文件存储路径
- 默认存储路径：`./upload`
- 可以在代码中修改存储路径

### 日志配置
- 日志文件路径：`./logs`
- 日志保留天数：30天
- 可以在`logger/logger.go`文件中修改日志配置

## 注意事项

1. **文件大小限制**：单个文件最大上传大小为100MB
2. **权限管理**：目前只有简单的管理员认证，建议在内部网络使用
3. **安全性**：建议在生产环境中配置HTTPS
4. **备份**：定期备份重要文件和日志
5. **性能**：建议根据服务器配置调整并发上传数

## 许可证

本项目采用 MIT 许可证，详情请查看 [LICENSE](https://github.com/Half-life7/gin_cloud_drive/blob/main/LICENSE) 文件。

MIT License 是一种宽松的开源许可证，允许您自由使用、修改和分发本项目的代码，无论是商业用途还是非商业用途，只要保留原始许可证和版权声明即可。

## 更新日志

### v1.0.0 (2025-12-01)
- ✅ 初始版本
- ✅ 实现文件上传下载功能
- ✅ 实现文件管理功能
- ✅ 实现系统状态监控
- ✅ 实现日志管理
- ✅ 实现暗色/亮色模式切换
- ✅ 实现拖拽上传
- ✅ 实现手动上传（文件缓冲区域）

## 快速开始

```bash
# 克隆项目
git clone https://github.com/Half-life7/gin_cloud_drive.git
cd gin_cloud_drive

# 安装依赖
go mod tidy

# 运行项目
go run main.go

# 访问系统
# 打开浏览器，访问 http://localhost:8080
```

## 贡献指南

1. **Fork 项目**
2. **创建特性分支** (`git checkout -b feature/AmazingFeature`)
3. **提交更改** (`git commit -m 'Add some AmazingFeature'`)
4. **推送到分支** (`git push origin feature/AmazingFeature`)
5. **创建 Pull Request**

### 贡献者行为准则

- 尊重他人，友好沟通
- 提交的代码需要有清晰的注释
- 确保代码通过编译和基本测试
- 提交的PR需要有清晰的描述

## 联系方式

- **GitHub**: [Half-life7](https://github.com/Half-life7)
- **E-mail**: tlingchenldb@gmail.com

## 致谢

- 感谢 [Gin](https://gin-gonic.com/) 框架团队
- 感谢 [gopsutil](https://github.com/shirou/gopsutil) 库的开发者
- 感谢 [Chart.js](https://www.chartjs.org/) 团队

---

**归2023级唐家鸿所有，未经许可不得复制、更改**