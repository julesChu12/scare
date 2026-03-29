package dto

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DateTime 时间类型，JSON 序列化为 RFC3339 秒级格式（无毫秒）
// 数据库存储为 DATETIME，与 time.Time 兼容
type DateTime struct {
	time.Time
}

// MarshalJSON 输出 RFC3339 秒级格式
func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format("2006-01-02T15:04:05Z07:00") + `"`), nil
}

// UnmarshalJSON 解析 RFC3339 格式
func (t *DateTime) UnmarshalJSON(data []byte) error {
	// 去掉引号
	s := string(data)
	if len(s) < 2 {
		return fmt.Errorf("invalid date format")
	}
	s = s[1 : len(s)-1]
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// Value 实现 driver.Valuer，支持写入数据库
func (t DateTime) Value() (driver.Value, error) {
	return t.Time, nil
}

// Scan 实现 sql.Scanner，支持从数据库读取
func (t *DateTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case []byte:
		parsed, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
	case string:
		parsed, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
	default:
		return fmt.Errorf("cannot scan %T into DateTime", value)
	}
	return nil
}
