package ahr999

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// AHR999公开数据API
	ahr999APIURL = "https://dncapi.flink1.com/api/v2/index/arh999?code=bitcoin&webp=1"
	// 缓存目录
	historyDir = "ahr999_history"
	// 月份文件模板
	monthFileTemplate = "%04d-%02d.json"
)

// Ahr999DataPoint 单条AHR999数据点
type Ahr999DataPoint struct {
	Date      string  `json:"date"`
	Timestamp int64   `json:"timestamp"`
	Ahr999    float64 `json:"ahr999"`
	BtcPrice  float64 `json:"btc_price"`
}

// GetAhr999 获取当前AHR999
func GetAhr999() (curr_btc_price, ahr999_value float64, err error) {
	// 缓存无则查API
	point, err := fetchAhr999FromAPI()
	if err != nil {
		return 0, 0, err
	}
	curr_btc_price = point.BtcPrice
	ahr999_value = point.Ahr999
	return curr_btc_price, ahr999_value, nil
}

// fetchAhr999FromAPI 获取API数据，并查询时间最近的一条
func fetchAhr999FromAPI() (*Ahr999DataPoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ahr999APIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求API失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	var response struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data [][]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %w", err)
	}
	if response.Code != 200 {
		return nil, fmt.Errorf("API返回错误: %s", response.Msg)
	}

	// 查询时间最近的一条（最大timestamp）
	var latest *Ahr999DataPoint
	for _, row := range response.Data {
		if len(row) >= 5 {
			timestamp, ok1 := row[0].(float64)
			ahr999, ok2 := row[1].(float64)
			btcPrice, ok3 := row[2].(float64)
			if ok1 && ok2 && ok3 {
				t := time.Unix(int64(timestamp), 0)
				date := t.Format("2006-01-02")
				point := Ahr999DataPoint{
					Date:      date,
					Timestamp: int64(timestamp),
					Ahr999:    ahr999,
					BtcPrice:  btcPrice,
				}
				if latest == nil || point.Timestamp > latest.Timestamp {
					p := point // create a copy to point to
					latest = &p
				}
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("API无有效数据")
	}
	return latest, nil
}
