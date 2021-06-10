package log

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"
	"unicode"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

type GormLogger struct {
	Location *time.Location
}

func (gl *GormLogger) Print(values ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			Error(fmt.Sprintf("gorm_logger::print::panic::%+v", err),
				String("err",
					fmt.Sprintf("%+v", err)),
				String("stack", string(debug.Stack())))
		}
	}()

	if len(values) <= 1 {
		return
	}

	var (
		sql             string
		formattedValues []string
		level           = fmt.Sprintf("%v", values[0])
		source          = fmt.Sprintf("%v", values[1])
	)

	if level == "sql" {
		if len(values) < 6 {
			return
		}

		// duration
		cost := float64(values[2].(time.Duration).Nanoseconds()/1e4) / 100.0
		// sql
		for _, value := range values[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			if indirectValue.IsValid() {
				value = indirectValue.Interface()
				if t, ok := value.(time.Time); ok {

					if gl.Location == nil {
						formattedValues = append(formattedValues, t.String())
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.In(gl.Location).Format("2006-01-02 15:04:05")))
					}
				} else if b, ok := value.([]byte); ok {
					if str := string(b); isPrintable(str) {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
					} else {
						formattedValues = append(formattedValues, "'<binary>'")
					}
				} else if r, ok := value.(driver.Valuer); ok {
					if value, err := r.Value(); err == nil && value != nil {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					} else {
						formattedValues = append(formattedValues, "NULL")
					}
				} else {
					switch value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
						formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
					default:
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				}
			} else {
				formattedValues = append(formattedValues, "NULL")
			}
		}

		// differentiate between $n placeholders or else treat like ?
		if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
			sql = values[3].(string)
			for index, value := range formattedValues {
				placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
				sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
			}
		} else {
			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		affectedRows := strconv.FormatInt(values[5].(int64), 10)
		Info("gorm-sql",
			Float64("cost", cost),
			String("affected", affectedRows),
			String("sql", sql),
			String("source", source))
	} else {
		Info("gorm-log",
			String("values", fmt.Sprintf("%v", values...)))
	}

	return
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
