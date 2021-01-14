package json

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	// "strings"
)

type JSONNullString struct {
	sql.NullString
}

func (ns *JSONNullString) MarshalJSON() ([]byte, error) {
	// formatted := fmt.Sprintf("\"%s\"", ns.String)
	// return []byte(formatted), nil
	if ns.Valid {
        return json.Marshal(ns.String)
    } else {
        // return json.Marshal(nil)
        return []byte(fmt.Sprintf("\"%s\"", "")), nil
    }
}

func (ns *JSONNullString) UnmarshalJSON(data []byte) error {
	// ns.String = strings.Trim(string(data), "\"") // Todo：此处有风险，"\"xxxx\"" 中间的两个"也会去掉
	// ns.Valid = false
	// return nil
	var s *string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    if s != nil {
		ns.String = *s
		ns.Valid = true
		// 这里是否兼容为""的情况？？
		if ns.String == "" {
			ns.Valid = false
		}
    } else {
        ns.Valid = false
    }
    return nil
}

// Value implements the driver Valuer interface.
func (ns JSONNullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// Scan implements the Scanner interface.
func (ns *JSONNullString) Scan(v interface{}) error {
	if v == nil {
		ns.String, ns.Valid = "", false
		return nil
	}

	ns.Valid = true
	return convertAssign(&ns.String, v)
}