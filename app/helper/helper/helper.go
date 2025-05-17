package helper

import (
	"time"
)

// 格式化日期时间
func LocalTimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 过滤map[string]interface{}类型的数据
func FilterMap(data map[string]interface{}, fields []string) map[string]interface{} {
	//参数过滤
	if len(fields) > 0 {
		for key, _ := range data {
			exist := false
			for _, field := range fields {
				if key == field {
					exist = true
					break
				}
			}
			if !exist {
				delete(data, key)
			}
		}
	}

	return data
}

// 合并多个map[string]interface{}
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
