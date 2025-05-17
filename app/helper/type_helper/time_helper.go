package type_helper

import (
	"database/sql/driver"
	"errors"
	"time"
)

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(timeFormat) + `"`), nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Time) Scan(value interface{}) error {
	tt, ok := value.(time.Time)
	if !ok {
		return errors.New("Time 类型有误")
	}

	*t = Time(tt)
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t Time) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}
