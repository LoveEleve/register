package main

import (
	"encoding/json"
	"fmt"
	"github.com/amarnathcjd/gogram/telegram"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var waitGroup sync.WaitGroup

type ActivationData struct {
	ActivationId string `json:"activationId"`
	PhoneNumber  string `json:"phoneNumber"`
}

func main() {
	waitGroup.Add(1)
	getMulNumber(1)
	waitGroup.Wait()
	//getCode("+959974600895")
	//Login()
	//getLoginCode()
}

// 获取多个号码
func getMulNumber(times int) {
	for count := 0; count < times; count++ {
		getSignleNumber()
	}
}

// 获取单个号码
func getSignleNumber() {

	// 基础API地址
	baseURL := "https://api.sms-activate.ae/stubs/handler_api.php"

	// 创建查询参数
	params := url.Values{}
	params.Set("api_key", "c6de5132c013098Ad74b4d4dc6b4b609") // 替换为你的实际API密钥
	params.Set("action", "getNumberV2")
	params.Set("service", "tg")
	params.Set("forward", "0")
	params.Set("operator", "") // 空值参数保留
	params.Set("ref", "")      // 空值参数保留
	params.Set("country", "187")
	params.Set("phoneException", "")
	params.Set("maxPrice", "")
	params.Set("useCashBack", "")
	params.Set("activationType", "")
	params.Set("language", "en-US")
	params.Set("userId", "") // 空值参数保留

	// 构造完整URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求对象
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	// 设置请求头（可选）
	req.Header.Set("User-Agent", "SMSActivateClient/1.0 (Go)")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	// 输出结果
	fmt.Printf("HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容:\n%s\n", string(body))

	var data ActivationData
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("JSON解析失败: " + err.Error())
	}
	// 处理电话号码
	if data.PhoneNumber != "" {
		formatted := "+" + data.PhoneNumber + "--" + data.ActivationId
		fmt.Println("处理后的号码:", formatted)
		//调用getCode()
		//go getCode(formatted)
		getCode(data)
	} else {
		fmt.Println("未找到有效号码")
	}
}

var client *telegram.Client

// login/register --> Login()函数会触发验证码的发送
func getCode(data ActivationData) {
	//waitGroup.Done()
	client, _ = telegram.NewClient(telegram.ClientConfig{ // Create a new client
		AppID:   1,
		AppHash: "b6b154c3707471f5339bd661645ed3d6",
		Session: data.PhoneNumber + "-session.dat",
	})
	client.Conn() // establish connection to the Telegram server
	client.Login(data.PhoneNumber, &telegram.LoginOptions{CodeCallback: func() (string, error) {
		fmt.Printf("Enter code: ")
		var codeInput string
		codeInput = getLoginCode(data)
		fmt.Scanln(&codeInput)
		return codeInput, nil
	}})
	fmt.Println(client.GetMe())
}

func Login(phone string) {
	client, _ := telegram.NewClient(telegram.ClientConfig{ // Create a new client
		AppID:   1,
		AppHash: "b6b154c3707471f5339bd661645ed3d6",
		Session: phone + "-session.dat",
	})
	client.Start()
}

/*
获取激活状态(验证码)
*/
func getLoginCode(data ActivationData) string {
	// 基础API地址
	baseURL := "https://api.sms-activate.ae/stubs/handler_api.php"

	// 创建查询参数
	params := url.Values{}
	params.Set("api_key", "c6de5132c013098Ad74b4d4dc6b4b609") // 需要替换为有效API密钥
	params.Set("action", "getActiveActivations")

	// 构造完整URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求对象
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return ""
	}

	// 设置请求头
	req.Header.Set("User-Agent", "SMSActivateClient/1.0 (Go)")
	req.Header.Set("Accept", "application/json") // 根据实际响应格式设置

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return ""
	}

	// 输出结果
	fmt.Printf("HTTP状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容:\n%s\n", string(body))
	return ""
}
