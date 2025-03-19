package filter

import (
	"bytes"
	"encoding/json"
	"github.com/beego/beego/v2/server/web/context"
	"io"
	"sl.framework.com/trace"
	"strconv"
)

//设置middleware 拦截请求并自动转换 string -> int64

func JsonMiddlewareConvert(ctx *context.Context) {

	//解析JSON
	trace.Debug("JsonMiddlewareConvert string to int64 begin url:%+v", ctx.Request.URL)
	var raw map[string]interface{}
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		if err.Error() != "EOF" {
			//请求体为空
			trace.Error("JsonMiddlewareConvert string to int64 url:%+v err:%+v", ctx.Request.URL, err)
		} else {
			//请求体stream已经被读取过
			trace.Notice("JsonMiddlewareConvert string to int64 url:%+v err:%+v", ctx.Request.URL, err)
		}
		return
	}
	trace.Debug("JsonMiddlewareConvert string to int64 url:%v raw:%+v", ctx.Request.URL, raw)
	//自动转换 string -> int64
	raw = convertStringToInt64(raw)
	//重新装箱 JSON
	modifiedBody, _ := json.Marshal(raw)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
}

func convertStringToInt64(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		switch v := value.(type) {
		case json.Number:
			//将int64转为string
			if _, err := v.Int64(); err == nil {
				data[key] = v.String()
				trace.Debug("convertStringToInt64 int64 to string key:%v,value:%v", key, value)
			}

		case string:
			// 尝试将字符串转换为 int64
			if num, err := strconv.ParseInt(v, 10, 64); err == nil {
				data[key] = num
				trace.Debug("convertStringToInt64 string to 64 key:%v,value:%v", key, value)
			}
		case map[string]interface{}:
			// 递归处理嵌套的 map
			data[key] = convertStringToInt64(v)
		case []interface{}:
			// 处理 slice，可能有字符串数字
			for i, elem := range v {
				if str, ok := elem.(string); ok {
					if num, err := strconv.ParseInt(str, 10, 64); err == nil {
						v[i] = num
					}
				} else if nestedMap, ok := elem.(map[string]interface{}); ok {
					v[i] = convertStringToInt64(nestedMap)
				}
			}
		}
	}
	return data
}
