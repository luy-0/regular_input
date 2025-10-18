package main

import (
	"log"
	"os"
	"os/signal"
	"regular_input/auto_buy"
	"regular_input/config"
	"regular_input/scheduler"
	"syscall"
	"time"
)

func main() {
	// 加载环境变量配置
	envConfig := config.LoadEnvConfig()

	// 检查是否存在 .env 文件
	if !config.CheckEnvFile() {
		log.Println("提示: 未找到 .env 文件，请参考 env.example 创建配置文件")
		config.PrintEnvSetupInstructions()
	} else {
		log.Println("提示: 已找到 .env 文件，使用环境变量配置")
		envConfig = config.LoadEnvConfig()
	}

	// 加载配置
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("警告: 加载配置文件失败: %v，使用默认配置", err)
		os.Exit(-1)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 打印环境变量状态
	envConfig.PrintEnvStatus(config.MessageConfig.PushMethod)

	// 验证环境变量（仅在非调试模式下验证）
	if !envConfig.Debug && !config.ParamsConfig.Debug {
		if err := envConfig.ValidateEnvConfig(config.GetPushMethods()); err != nil {
			log.Printf("环境变量验证失败: %v", err)
			missingVars := envConfig.GetMissingEnvVars(config.GetPushMethods())
			if len(missingVars) > 0 {
				log.Printf("缺少的环境变量: %v", missingVars)
			}
			log.Fatal("请设置必需的环境变量后重试")
		}
	} else {
		log.Println("调试模式：跳过环境变量验证")
	}

	log.Printf("已加载配置: %s", config.TaskConfig.Name)
	log.Printf("定时表达式: %s", config.TaskConfig.Schedule)
	log.Printf("基础金额: %.2f USDT", config.ParamsConfig.BaseAmount)
	log.Printf("使用 AHR999: %v", config.ParamsConfig.UseAhr999)
	log.Printf("调试模式: %v", config.ParamsConfig.Debug)

	// 创建定投任务
	apiKey := envConfig.BinanceAPIKey
	secretKey := envConfig.BinanceSecretKey
	proxyUrl := envConfig.HTTPSProxy

	autoBuyConfig := &auto_buy.Config{
		Name:             config.TaskConfig.Name,
		BaseAmount:       config.ParamsConfig.BaseAmount,
		UseAhr999:        config.ParamsConfig.UseAhr999,
		Ahr999TimerTable: config.ParamsConfig.Ahr999TimerTable,
		Symbol:           "BTCUSDT",
		Debug:            config.ParamsConfig.Debug,
		PushMethods:      config.GetPushMethods(),
	}

	autoBuyTask := auto_buy.NewAutoBuyTask(autoBuyConfig, apiKey, secretKey, proxyUrl)

	// 创建调度器
	sched := scheduler.NewScheduler()

	// 添加定投任务到调度器
	if err := sched.AddTask(
		config.TaskConfig.Name,
		config.TaskConfig.Schedule,
		autoBuyTask,
	); err != nil {
		log.Fatalf("添加任务失败: %v", err)
	}

	// 启动调度器
	sched.Start()
	log.Printf("定时任务调度器已启动，任务将在 %s 执行", config.TaskConfig.Schedule)

	// 如果是调试模式，立即执行一次任务
	if config.ParamsConfig.ExecuteNow {
		log.Println("调试模式：立即执行一次任务...")
		go func() {
			time.Sleep(2 * time.Second) // 等待2秒让调度器完全启动
			if err := sched.ExecuteTaskNow(config.TaskConfig.Name); err != nil {
				log.Printf("执行任务失败: %v", err)
			}
		}()
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("按 Ctrl+C 停止...")
	<-sigChan

	// 优雅停止
	log.Println("正在停止定时任务调度器...")
	sched.Stop()
	log.Println("定时任务调度器已停止")
}
