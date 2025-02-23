package main

import (
    "bytes"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "net/http"
)

func main() {
    router := gin.Default()

    // 定义一个 POST 路由，接收用户消息并调用外部 API
    router.POST("/chat", func(c *gin.Context) {

        // 获取请求参数
        var request struct {
            Content string `json:"content"`
        }
        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 构造外部 API 请求
        apiUrl := "https://api.siliconflow.cn/v1/chat/completions"
        payload := map[string]interface{}{
            "model": "deepseek-ai/DeepSeek-R1-Distill-Llama-8B",
            "messages": []map[string]string{
                {"role": "user", "content": request.Content},
            },
        }

        // 设置请求头和请求体
        jsonPayload, _ := json.Marshal(payload)
        req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonPayload))
        req.Header.Set("Authorization", "Bearer sk-mglfvvigpgmtelqwwjdjsblqkhplzciabixdgzgunoeglahp")
        req.Header.Set("Content-Type", "application/json")

        // 发起请求
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
            return
        }
        defer resp.Body.Close()

        // 读取响应内容
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != http.StatusOK {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "External API returned an error"})
            return
        }

        // 返回外部 API 的响应
        c.Data(resp.StatusCode, "application/json", body)
    })

    // 启动服务器
    router.Run(":8080")
}