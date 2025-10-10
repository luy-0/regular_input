package main

import (
	"os"
	"task_scheduler/helper"
)

func main() {
	// 检查必要的环境变量
	if os.Getenv("TELEGRAM_BOT_TOKEN") == "" && 
	   os.Getenv("LARK_TOKEN") == "" && 
	   os.Getenv("WECHAT_SEND_KEY") == "" {
		println("提示: 请设置环境变量 TELEGRAM_BOT_TOKEN, LARK_TOKEN 或 WECHAT_SEND_KEY")
		println("示例:")
		println("  export TELEGRAM_BOT_TOKEN=your_token")
		println("  export TELEGRAM_CHAT_ID=your_chat_id")
		println("  export LARK_TOKEN=your_lark_token")
		println("  export WECHAT_SEND_KEY=your_wechat_key")
		println()
	}

	// 运行推送器测试
	helper.TestAllPushers()
}

