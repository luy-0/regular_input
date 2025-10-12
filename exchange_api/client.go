package exchange_api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

// Client CCXT交易所客户端
type Client struct {
	spotClient *binance.Client
	apiKey     string
	secretKey  string
	proxyUrl   string
}

// NewClient 创建新的CCXT客户端
func NewClient(apiKey, secretKey, proxyUrl string) *Client {
	// 只有在明确设置了代理时才使用代理
	if proxyUrl == "" {
		proxyUrl = os.Getenv("HTTPS_PROXY")
	}

	if apiKey == "" || secretKey == "" {
		log.Println("apiKey or secretKey is empty")
	}

	var cli *Client
	if proxyUrl != "" {
		// 使用代理客户端
		log.Printf("使用代理连接: %s", proxyUrl)
		cli = &Client{
			spotClient: binance.NewProxiedClient(apiKey, secretKey, proxyUrl),
			apiKey:     apiKey,
			secretKey:  secretKey,
			proxyUrl:   proxyUrl,
		}
	} else {
		// 使用直连客户端
		log.Println("使用直连模式（无代理）")
		cli = &Client{
			spotClient: binance.NewClient(apiKey, secretKey),
			apiKey:     apiKey,
			secretKey:  secretKey,
			proxyUrl:   "",
		}
	}
	cli.Ping(context.Background())
	return cli
}

// NewClientWithoutAuth 创建无需认证的客户端（仅用于公开接口）
func NewClientWithoutAuth(proxyUrl string) *Client {
	// 只有在明确设置了代理时才使用代理
	if proxyUrl == "" {
		proxyUrl = os.Getenv("HTTPS_PROXY")
	}
	ret := &Client{}
	if proxyUrl != "" {
		// 使用代理客户端
		log.Printf("使用代理连接（无认证）: %s", proxyUrl)
		ret.spotClient = binance.NewProxiedClient("", "", proxyUrl)
		ret.proxyUrl = proxyUrl
	} else {
		// 使用直连客户端
		log.Println("使用直连模式（无认证，无代理）")
		ret.spotClient = binance.NewClient("", "")
		ret.proxyUrl = ""
	}
	ret.Ping(context.Background())
	return ret
}

// Ping 测试交易所连通性
func (c *Client) Ping(ctx context.Context) error {
	// 使用Ping接口测试连通性
	err := c.spotClient.NewPingService().Do(ctx)
	if err != nil {
		return fmt.Errorf("ping交易所失败: %w", err)
	}
	prices, err := c.GetBTCPrice(ctx)
	if err != nil {
		return fmt.Errorf("获取BTC价格失败: %w", err)
	}
	log.Println("🎉交易所连通成功，当前BTC价格:", prices)
	return nil
}

// GetLatestPrice 获取最新价格
func (c *Client) GetLatestPrice(ctx context.Context, symbol string) (*binance.SymbolPrice, error) {
	// 获取单个交易对的最新价格
	price, err := c.spotClient.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取%s价格失败: %w", symbol, err)
	}

	if len(price) == 0 {
		return nil, fmt.Errorf("未找到%s的价格数据", symbol)
	}

	return price[0], nil
}

// GetLatestPrices 获取多个交易对的最新价格
func (c *Client) GetLatestPrices(ctx context.Context) ([]*binance.SymbolPrice, error) {
	// 获取所有交易对的最新价格
	prices, err := c.spotClient.NewListPricesService().Symbol("BTCUSDT").Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取价格列表失败: %w", err)
	}

	return prices, nil
}

// GetBTCPrice 获取比特币最新价格
func (c *Client) GetBTCPrice(ctx context.Context) (float64, error) {
	price, err := c.GetLatestPrice(ctx, "BTCUSDT")
	if err != nil {
		return 0, err
	}

	// 将字符串价格转换为float64
	priceFloat, err := parseFloat(price.Price)
	if err != nil {
		return 0, fmt.Errorf("解析价格失败: %w", err)
	}

	return priceFloat, nil
}

// GetETHPrice 获取以太坊最新价格
func (c *Client) GetETHPrice(ctx context.Context) (float64, error) {
	price, err := c.GetLatestPrice(ctx, "ETHUSDT")
	if err != nil {
		return 0, err
	}

	// 将字符串价格转换为float64
	priceFloat, err := parseFloat(price.Price)
	if err != nil {
		return 0, fmt.Errorf("解析价格失败: %w", err)
	}

	return priceFloat, nil
}

// GetServerTime 获取服务器时间
func (c *Client) GetServerTime(ctx context.Context) (time.Time, error) {
	serverTime, err := c.spotClient.NewServerTimeService().Do(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("获取服务器时间失败: %w", err)
	}

	return time.Unix(serverTime/1000, (serverTime%1000)*1000000), nil
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
	// 1. 测试ping连通性
	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("ping测试失败: %w", err)
	}

	// 2. 测试获取服务器时间
	if _, err := c.GetServerTime(ctx); err != nil {
		return fmt.Errorf("获取服务器时间失败: %w", err)
	}

	// 3. 测试获取BTC价格
	if _, err := c.GetBTCPrice(ctx); err != nil {
		return fmt.Errorf("获取BTC价格失败: %w", err)
	}

	return nil
}

// GetKlines 获取K线数据
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*binance.Kline, error) {
	klines, err := c.spotClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取%s K线数据失败: %w", symbol, err)
	}

	return klines, nil
}

// GetKlinesWithTimeRange 获取指定时间范围的K线数据
func (c *Client) GetKlinesWithTimeRange(ctx context.Context, symbol string, interval string, startTime, endTime time.Time) ([]*binance.Kline, error) {
	klines, err := c.spotClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		StartTime(startTime.UnixMilli()).
		EndTime(endTime.UnixMilli()).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取%s K线数据失败: %w", symbol, err)
	}

	return klines, nil
}

// GetBTCHistoryPrices 获取比特币历史价格数据（用于AHR999计算）
func (c *Client) GetBTCHistoryPrices(ctx context.Context, days int) ([][]float64, error) {
	// 计算需要的K线数量（每天1根K线）
	limit := days

	// 获取日K线数据
	klines, err := c.GetKlines(ctx, "BTCUSDT", "1d", limit)
	if err != nil {
		return nil, fmt.Errorf("获取BTC历史K线数据失败: %w", err)
	}

	// 转换为所需的格式 [][]float64，每个元素为 [timestamp, price]
	var prices [][]float64
	for _, kline := range klines {
		// 将时间戳转换为毫秒
		timestamp := float64(kline.OpenTime)

		// 使用收盘价作为当日价格
		price, err := parseFloat(kline.Close)
		if err != nil {
			continue // 跳过无效价格
		}

		prices = append(prices, []float64{timestamp, price})
	}

	return prices, nil
}

// GetBTCHistoryPricesForDate 获取指定日期附近的比特币价格数据
func (c *Client) GetBTCHistoryPricesForDate(ctx context.Context, targetDate time.Time, daysBefore int) ([][]float64, error) {
	// 计算时间范围
	endTime := targetDate.Add(24 * time.Hour)          // 目标日期的下一天
	startTime := targetDate.AddDate(0, 0, -daysBefore) // daysBefore天前

	// 获取指定时间范围的K线数据
	klines, err := c.GetKlinesWithTimeRange(ctx, "BTCUSDT", "1d", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("获取BTC历史K线数据失败: %w", err)
	}

	// 转换为所需的格式 [][]float64，每个元素为 [timestamp, price]
	var prices [][]float64
	for _, kline := range klines {
		// 将时间戳转换为毫秒
		timestamp := float64(kline.OpenTime)

		// 使用收盘价作为当日价格
		price, err := parseFloat(kline.Close)
		if err != nil {
			continue // 跳过无效价格
		}

		prices = append(prices, []float64{timestamp, price})
	}

	return prices, nil
}

// GetBTCPriceAtDate 获取指定日期的比特币价格
func (c *Client) GetBTCPriceAtDate(ctx context.Context, targetDate time.Time) (float64, error) {
	// 获取目标日期附近的价格数据（前后各1天）
	endTime := targetDate.Add(24 * time.Hour)
	startTime := targetDate.Add(-24 * time.Hour)

	klines, err := c.GetKlinesWithTimeRange(ctx, "BTCUSDT", "1d", startTime, endTime)
	if err != nil {
		return 0, fmt.Errorf("获取BTC历史价格失败: %w", err)
	}

	// 找到最接近目标日期的价格
	var closestPrice float64
	var minDiff time.Duration = 24 * time.Hour

	for _, kline := range klines {
		klineTime := time.Unix(kline.OpenTime/1000, 0)
		diff := absDuration(targetDate.Sub(klineTime))

		if diff < minDiff {
			price, err := parseFloat(kline.Close)
			if err != nil {
				continue
			}
			closestPrice = price
			minDiff = diff
		}
	}

	if closestPrice == 0 {
		return 0, fmt.Errorf("未找到%s附近的价格数据", targetDate.Format("2006-01-02"))
	}

	return closestPrice, nil
}

// absDuration 计算时间差的绝对值
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// parseFloat 辅助函数：将字符串转换为float64
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		fmt.Println("解析价格失败: ", err)
		return 0, err
	}
	return f, nil
}

// jsonAnything 辅助函数：将任何类型转换按照美化后的json格式输出
func jsonAnything(v any) string {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "json marshal error"
	}
	return string(json)
}

// 获取账户金额信息
func (c *Client) GetAccountBalance(ctx context.Context) string {
	// 传入 omitZeroBalances == true 来过滤掉零余额的资产
	balance, err := c.spotClient.NewGetAccountService().OmitZeroBalances(true).Do(ctx)
	if err != nil {
		return fmt.Sprintf("获取账户金额信息失败: %v", err)
	}
	return jsonAnything(balance)
}

// 获取账户的BTC余额
func (c *Client) GetBTCBalance(ctx context.Context) string {
	account, err := c.spotClient.NewGetAccountService().OmitZeroBalances(true).Do(ctx)
	if err != nil {
		return fmt.Sprintf("获取账户失败: %v", err)
	}
	balance := account.Balances
	for _, balance := range balance {
		if balance.Asset == "BTC" {
			return balance.Free
		}
	}
	return "0"
}

// 依照传入的 Symbol 和 Amount 按照市价购买指定数量的币
// 参数: symbol: 币种名称(BTCUSDT), amount: 购买数量(USDT)
func (c *Client) BuyCoinByMarketPrice(ctx context.Context, symbol string, amount float64) (string, error) {
	order, err := c.spotClient.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(strconv.FormatFloat(amount, 'f', -1, 64)).Do(ctx)
	if err != nil {
		return fmt.Sprintf("购买%s失败: %v", symbol, err), err
	}
	return jsonAnything(order), nil
}

// 依照传入的 Symbol 和 Amount 按照最优价购买指定数量的币
// 参数: symbol: 币种名称(BTCUSDT), amount: 购买金额(USDT)
// 首先获取当前订单盘口，然后根据盘口价格计算最优价
func (c *Client) BuyCoinByBestPrice(ctx context.Context, symbol string, amount float64) (string, error) {
	bestSellPrice, _, err := c.GetBestPrice(ctx, symbol)
	if err != nil {
		return fmt.Sprintf("获取%s订单盘口失败: %v", symbol, err), err
	}
	acount := amount / bestSellPrice
	acount = math.Round(acount*100000) / 100000
	fmt.Printf("购买数量: %.5f\n", acount)
	order, err := c.spotClient.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Price(strconv.FormatFloat(bestSellPrice, 'f', -1, 64)).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(strconv.FormatFloat(acount, 'f', -1, 64)).
		Do(ctx)
	if err != nil {
		return fmt.Sprintf("购买%s失败: %v", symbol, err), err
	}
	return jsonAnything(order), nil
}

// 获取当前订单盘口
func (c *Client) GetBestPrice(ctx context.Context, symbol string) (bestSellPrice, bestBuyPrice float64, err error) {
	orderBook, err := c.spotClient.NewListBookTickersService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("获取%s订单盘口失败: %v", symbol, err)
	}
	bestSellPrice, err = parseFloat(orderBook[0].AskPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("获取%s订单盘口失败: %v", symbol, err)
	}
	bestBuyPrice, err = parseFloat(orderBook[0].BidPrice)
	if err != nil {
		return 0, 0, fmt.Errorf("获取%s订单盘口失败: %v", symbol, err)
	}
	return
}
