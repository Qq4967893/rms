package formatter

import (
	"bytes"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	DateTime = "2006-01-02 15:04:05"
	DateOnly = "2006-01-02"
	TimeOnly = "15:04:05"
)

func Datetime(t time.Time) string {
	return t.Format(DateTime)
}

func Amount(d decimal.Decimal) string {
	buf := &bytes.Buffer{}
	if d.IsNegative() {
		buf.Write([]byte{'-'})
		d = d.Neg()
	}
	comma := []byte{','}
	parts := strings.Split(d.StringFixedBank(2), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)
	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func Bool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
