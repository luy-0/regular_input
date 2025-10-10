package helper

import (
	"fmt"
	"os"
)

// TestAllPushers 测试所有推送器
func TestAllPushers() {
	fmt.Println("=== 开始测试推送器 ===\n")

	results := make(map[string]bool)

	// 测试 Telegram 推送器
	// 直接设置 HTTP/HTTPS 代理为本地 10809 端口，便于测试 Telegram/Lark 等推送
	// proxyAddr := "http://127.0.0.1:7890"
	// _ = proxyAddr
	// // 设置环境变量
	// setProxy := func() {
	// 	// 直接用 os.Setenv，不在此处 import "os"，但 main.go 已经 import，可以生效。
	// 	// 推荐将该代码移动到 main() 入口处做全局设置。如果只是临时测试推送可用此法:
	// 	fmt.Println("⚠️  已为测试推送设置本地 HTTP_PROXY/HTTPS_PROXY = " + proxyAddr)
	// 	os.Setenv("HTTP_PROXY", proxyAddr)
	// 	os.Setenv("HTTPS_PROXY", proxyAddr)
	// }
	// setProxy()
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramToken != "" && telegramChatID != "" {
		fmt.Println("📱 测试 Telegram 推送器...")
		telegramBot := NewTelegramBot(telegramToken, telegramChatID)

		success := telegramBot.TestPush()
		results["Telegram"] = success
		if success {
			fmt.Println("✅ Telegram 推送测试成功")
		} else {
			fmt.Println("❌ Telegram 推送测试失败")
		}
	} else {
		fmt.Println("⚠️  跳过 Telegram 测试 (未配置 TELEGRAM_BOT_TOKEN 或 TELEGRAM_CHAT_ID)")
		results["Telegram"] = false
	}

	fmt.Println()

	// 测试 Lark (飞书) 推送器
	larkToken := os.Getenv("LARK_TOKEN")
	if larkToken != "" {
		fmt.Println("🔔 测试 Lark (飞书) 推送器...")
		larkBot := NewLarkBot(larkToken)

		success := larkBot.TestPush()
		results["Lark"] = success
		if success {
			fmt.Println("✅ Lark 推送测试成功")
		} else {
			fmt.Println("❌ Lark 推送测试失败")
		}
	} else {
		fmt.Println("⚠️  跳过 Lark 测试 (未配置 LARK_TOKEN)")
		results["Lark"] = false
	}

	fmt.Println()

	// 测试微信推送器
	wechatSendKey := os.Getenv("WECHAT_SEND_KEY")
	if wechatSendKey != "" {
		fmt.Println("💬 测试微信推送器...")
		wechatPusher := NewWeChatPusher(wechatSendKey)

		success := wechatPusher.TestPush()
		// success := true
		results["WeChat"] = success
		if success {
			fmt.Println("✅ 微信推送测试成功")
		} else {
			fmt.Println("❌ 微信推送测试失败")
		}
	} else {
		fmt.Println("⚠️  跳过微信测试 (未配置 WECHAT_SEND_KEY)")
		results["WeChat"] = false
	}

	// 打印总结
	fmt.Println("\n=== 测试结果汇总 ===")
	successCount := 0
	totalCount := 0
	for pusher, success := range results {
		totalCount++
		status := "❌ 失败"
		if success {
			status = "✅ 成功"
			successCount++
		}
		fmt.Printf("%s: %s\n", pusher, status)
	}
	fmt.Printf("\n总计: %d/%d 推送器测试通过\n", successCount, totalCount)
}
