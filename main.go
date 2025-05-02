package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

type WebToCodeRequest struct {
	URL         string `json:"url" description:"要解析的网页 URL" required:"true"`
	Language    string `json:"language" description:"目标编程语言，如 java、python、go 等" required:"true"`
	ExtraPrompt string `json:"extraPrompt" description:"额外的提示信息，可选" required:"false"`
}

func main() {
	// 创建 SSE 传输服务器
	transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("创建传输服务器失败: %v", err)
	}

	// 初始化 MCP 服务器
	mcpServer, err := server.NewServer(transportServer)
	if err != nil {
		log.Fatalf("创建 MCP 服务器失败: %v", err)
	}

	// 注册网页到代码转换工具
	tool, err := protocol.NewTool("web_to_code", "解析网页内容并生成相应的代码及注释", WebToCodeRequest{})
	if err != nil {
		log.Fatalf("创建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(tool, handleWebToCodeRequest)

	// 启动服务器
	log.Println("服务器运行在 http://127.0.0.1:8080")
	if err = mcpServer.Run(); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}

func handleWebToCodeRequest(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var webRequest WebToCodeRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &webRequest); err != nil {
		return nil, err
	}

	// 1. 获取网页内容
	htmlContent, err := fetchWebContent(webRequest.URL)
	if err != nil {
		return nil, fmt.Errorf("获取网页内容失败: %v", err)
	}

	// 2. 提取网页中的主要文本内容
	textContent, err := extractTextFromHTML(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("提取文本内容失败: %v", err)
	}

	// 3. 添加提示词和要求
	prompt := fmt.Sprintf("请根据以下网页内容生成%s代码:\n\n", webRequest.Language)
	requirements := "\n\n## 代码生成要求\n\n" +
		"1. 代码需要面向初学者，添加详细注释解释每个关键部分\n" +
		"2. 使用清晰的命名和良好的代码结构\n" +
		"3. 包含必要的错误处理和边界条件检查\n"

	if webRequest.ExtraPrompt != "" {
		requirements += fmt.Sprintf("4. 额外要求: %s\n", webRequest.ExtraPrompt)
	}

	fullContent := prompt + textContent + requirements

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: fullContent,
			},
		},
	}, nil
}

// 获取网页内容
func fetchWebContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 从 HTML 中提取文本内容
func extractTextFromHTML(htmlContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// 移除 script 和 style 标签
	doc.Find("script, style").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// 获取主要内容区域（这里需要根据目标网站进行调整）
	// 例如，可以尝试定位 article, main, .content 等常见内容区域标签
	var content string
	mainContent := doc.Find("article, main, .content, #content").First()
	if mainContent.Length() > 0 {
		content = mainContent.Text()
	} else {
		content = doc.Find("body").Text()
	}

	// 清理文本
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\n\n\n", "\n\n")

	return content, nil
}
