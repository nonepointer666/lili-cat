package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Choices struct {
	index         int
	message       Message
	finish_reason string
}

type Usage struct {
	prompt_tokens     int
	completion_tokens int
	total_tokens      int
}

type Message struct {
	role              string
	content           string
	reasoning_content string
}

type OpenAIText struct {
	id      string
	object  string
	created int
	model   string
	choices []Choices
	usage   Usage
}

type Error struct {
	code    string
	message string
}

// 将 map[string]interface{} 转换为 OpenAIText
func mapToOpenAIText(data map[string]interface{}) (OpenAIText, error) {
	var result OpenAIText

	// 使用反射遍历 map
	for key, value := range data {
		switch key {
		case "id":
			if id, ok := value.(string); ok {
				result.id = id
			} else {
				return result, fmt.Errorf("invalid type for 'id'")
			}
		case "object":
			if object, ok := value.(string); ok {
				result.object = object
			} else {
				return result, fmt.Errorf("invalid type for 'object'")
			}
		case "created":
			if created, ok := value.(float64); ok {
				result.created = int(created)
			} else {
				return result, fmt.Errorf("invalid type for 'created'")
			}
		case "model":
			if model, ok := value.(string); ok {
				result.model = model
			} else {
				return result, fmt.Errorf("invalid type for 'model'")
			}
		case "choices":
			if choices, ok := value.([]interface{}); ok {
				result.choices = make([]Choices, len(choices))
				for i, choice := range choices {
					if choiceMap, ok := choice.(map[string]interface{}); ok {
						result.choices[i] = mapToChoices(choiceMap)
					} else {
						return result, fmt.Errorf("invalid type for 'choices'")
					}
				}
			} else {
				return result, fmt.Errorf("invalid type for 'choices'")
			}
		case "usage":
			if usageMap, ok := value.(map[string]interface{}); ok {
				result.usage = mapToUsage(usageMap)
			} else {
				return result, fmt.Errorf("invalid type for 'usage'")
			}
		default:
			fmt.Printf("Unknown key: %s\n", key)
		}
	}

	return result, nil
}

// 将 map 转换为 Choices
func mapToChoices(data map[string]interface{}) Choices {
	var result Choices

	if index, ok := data["index"].(float64); ok {
		result.index = int(index)
	}

	if messageMap, ok := data["message"].(map[string]interface{}); ok {
		result.message = mapToMessage(messageMap)
	}

	if finishReason, ok := data["finish_reason"].(string); ok {
		result.finish_reason = finishReason
	}

	return result
}

// 将 map 转换为 Usage
func mapToUsage(data map[string]interface{}) Usage {
	var result Usage

	if promptTokens, ok := data["prompt_tokens"].(float64); ok {
		result.prompt_tokens = int(promptTokens)
	}

	if completionTokens, ok := data["completion_tokens"].(float64); ok {
		result.completion_tokens = int(completionTokens)
	}

	if totalTokens, ok := data["total_tokens"].(float64); ok {
		result.total_tokens = int(totalTokens)
	}

	return result
}

// 将 map 转换为 Message
func mapToMessage(data map[string]interface{}) Message {
	var result Message

	if role, ok := data["role"].(string); ok {
		result.role = role
	}

	if content, ok := data["content"].(string); ok {
		result.content = content
	}

	if reasoningContent, ok := data["reasoning_content"].(string); ok {
		result.reasoning_content = reasoningContent
	}

	return result
}
func getOpenAITextModel() OpenAIText {
	// 初始化 Message
	message := Message{
		role:              "assistant",
		content:           "Hello, how can I help you today?",
		reasoning_content: "This is a friendly greeting message.",
	}

	// 初始化 Choices
	choices := []Choices{
		{
			index:         0,
			message:       message,
			finish_reason: "stop",
		},
	}

	// 初始化 Usage
	usage := Usage{
		prompt_tokens:     10,
		completion_tokens: 20,
		total_tokens:      30,
	}

	// 初始化 OpenAIText
	openAIText := OpenAIText{
		id:      "",
		object:  "",
		created: 0, // 当前时间戳
		model:   "",
		choices: choices,
		usage:   usage,
	}

	return openAIText
}

func requestChatAPI(content string) (*OpenAIText, error) {
	// 构造请求数据（假设接口需要一个 "prompt" 参数）
	// 构造外部 API 请求
	apiUrl := "https://api.siliconflow.cn/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "deepseek-ai/DeepSeek-R1-Distill-Llama-8B",
		"messages": []map[string]string{
			{"role": "user", "content": content},
		},
	}

	// 设置请求头和请求体
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer sk-mglfvvigpgmtelqwwjdjsblqkhplzciabixdgzgunoeglahp")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	fmt.Println(body)
	// 解析响应为 OpenAIText 结构
	var result OpenAIText
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	fmt.Println(body)

	return &result, nil
}

// 请求 Chat API
func requestChatAPI2(content string) (map[string]interface{}, error) {
	client := &http.Client{}

	// 构造外部 API 请求
	apiUrl := "https://api.siliconflow.cn/v1/chat/completions"
	payload := map[string]interface{}{
		"model": "deepseek-ai/DeepSeek-R1-Distill-Llama-8B",
		"messages": []map[string]string{
			{"role": "user", "content": content},
		},
	}

	// 设置请求头和请求体
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer sk-mglfvvigpgmtelqwwjdjsblqkhplzciabixdgzgunoeglahp")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return response, nil
}

func main() {

	router := gin.Default()

	// 设置信任所有代理
	router.SetTrustedProxies(nil)

	// 定义一个 POST 路由，接收用户消息并调用外部 API
	router.POST("/chat", func(c *gin.Context) {

		// 获取请求参数
		var request struct {
			Content string `json:"content"`
		}
		resp, err := requestChatAPI2(request.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// 返回外部 API 的响应
		c.JSON(http.StatusOK, resp)
	})

	// 定义一个 POST 路由，接收用户消息并调用外部 API
	router.POST("/sse", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		// 获取请求参数
		var request struct {
			Content string `json:"content"`
		}
		resp, err := requestChatAPI2(request.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 使用反射遍历 map
		text, err := mapToOpenAIText(resp)
		if len(text.choices) > 0 {
			// 获取第一个 choice 的 message.content
			messageContent := text.choices[0].message.content
			fmt.Printf("---------------------------")

			fmt.Println(messageContent)
			fmt.Printf("---------------------------")

			// 模拟数据推送
			for index, char := range messageContent {
				// 格式化为 SSE 数据
				data := fmt.Sprintf("data: %c\n\n", char)
				_, err := c.Writer.WriteString(data)
				if err != nil {
					// 如果发生错误，退出循环
					fmt.Println("Error writing to client:", index, err)
					return
				}
				fmt.Println("sse", index, char)

				c.Writer.Flush() // 立即推送数据
			}
		} else {
			fmt.Println("No choices available.")
		}

	})

	// 启动服务器
	router.Run(":8080")
}
