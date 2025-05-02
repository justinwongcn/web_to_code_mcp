# Web to Code MCP

## 项目简介

这是一个基于Go语言的网页转代码工具，通过MCP服务提供网页内容解析和代码生成功能。

## 核心功能

- 解析网页内容并提取主要文本
- 根据指定编程语言生成代码
- 自动添加详细注释和错误处理
- 支持通过MCP服务远程调用

## 安装指南

1. 确保已安装Go 1.24或更高版本
2. 克隆项目仓库
3. 安装依赖:
```bash
go mod download
```

## 使用方法

1. 启动服务:
```bash
go run main.go
```
2. 通过MCP客户端调用`web_to_code`工具
3. 提供网页URL和目标编程语言

## MCP 配置

```json
{
   "mcpServers": { 
     "web_to_code": { 
       "url": "http://127.0.0.1:8080/sse" 
     } 
   } 
}
```

## 依赖

- github.com/PuerkitoBio/goquery
- github.com/ThinkInAIXYZ/go-mcp

## 许可证

MIT License (见LICENSE文件)