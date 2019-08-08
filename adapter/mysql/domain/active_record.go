package domain

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/visionom/v-gae/adapter/mysql/errors"
)

type ActiveRecordList []*ActiveRecord

func (rds ActiveRecordList) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "[")
			i := 0

			for _, rd := range rds {
				fmt.Fprintf(&buf, "{")
				j := 0
				for k, v := range rd.rows {
					fmt.Fprintf(&buf, "\n  %s:%s", k, string((v)[:]))
					if j++; j < len(rd.rows) {
						fmt.Fprintf(&buf, ",")
					}
				}
				fmt.Fprintf(&buf, "\n}")
				if i++; i < len(rds) {
					fmt.Fprintf(&buf, ",")
				}
			}

			io.WriteString(s, buf.String())
			return
		}

		fallthrough
	case 's', 'q':
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "[")
		i := 0
		for _, rd := range rds {
			fmt.Fprintf(&buf, "{")
			j := 0
			for k, v := range rd.rows {
				fmt.Fprintf(&buf, "%s:%s", k, string((v)[:]))
				if j++; j < len(rd.rows) {
					fmt.Fprintf(&buf, ",")
				}
			}
			fmt.Fprintf(&buf, "}")
			if i++; i < len(rds) {
				fmt.Fprintf(&buf, ",")
			}
		}
		fmt.Fprintf(&buf, "]")
		io.WriteString(s, buf.String())
	}
}

type ActiveRecord struct {
	cols []string
	rows map[string][]byte
}

func NewRd() *ActiveRecord {
	rows := make(map[string][]byte)
	return &ActiveRecord{[]string{}, rows}
}

func GetActiveRecordList(rows *sql.Rows) (ActiveRecordList, error) {
	rds := make([]*ActiveRecord, 0)
	if rows == nil {
		return []*ActiveRecord{},
			errors.New("rows is empty")
	}

	for rows.Next() {
		rd, err := parse(rows)
		if err != nil {
			return nil, err
		}
		rds = append(rds, rd)
	}
	return rds, nil
}

func GetActiveRecord(rows *sql.Rows) (*ActiveRecord, error) {
	if rows.Next() {
		return parse(rows)
	}
	return &ActiveRecord{[]string{}, map[string][]byte{}}, nil
}

func parse(rows *sql.Rows) (*ActiveRecord, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.New("get columns from Rows fails")
	}

	lenCN := len(cols)
	raw := newRaw(lenCN)
	err = rows.Scan(raw...)
	if err != nil {
		return nil,
			errors.New("scan data from rows fails")
	}

	rd := &ActiveRecord{
		cols: cols,
		rows: make(map[string][]byte),
	}
	for i := 0; i < lenCN; i++ {
		var buf bytes.Buffer
		_, err = buf.Write(*raw[i].(*sql.RawBytes))
		if err != nil {
			return nil, errors.New("write data to buff fails")
		}
		rd.rows[strings.ToLower(cols[i])] = buf.Bytes()
	}
	return rd, nil
}

func newRaw(lenCN int) []interface{} {
	raw := make([]interface{}, lenCN)
	for i := 0; i < lenCN; i++ {
		raw[i] = new(sql.RawBytes)
	}
	return raw
}

func (rd *ActiveRecord) GetCols() []string {
	return rd.cols
}

func (rd *ActiveRecord) GetRows() []string {
	return rd.cols
}

func (rd *ActiveRecord) Get(colName string) (string,
	error) {
	if colName == "" {
		return "", errors.New("col name is empty")
	}

	if value, ok := rd.rows[strings.ToLower(colName)]; ok {
		result := string((value)[:])
		return result, nil
	} else {
		return "", errors.Newf("can't find '%s' in this record", colName)
	}
}

func (rd *ActiveRecord) GetString(colName string) (string, error) {
	return rd.Get(colName)
}

func (rd *ActiveRecord) GetInt(colName string) (int, error) {
	value, err := rd.Get(colName)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Newf("convert %s to int fails, value is %s", colName, value)
	}
	return v, nil
}

func (rd *ActiveRecord) GetInt64(colName string) (int64, error) {
	value, err := rd.Get(colName)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Newf("convert %s to int64 fails, value is %s", colName, value)
	}
	return v, nil
}

func (rd *ActiveRecord) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			var buf bytes.Buffer
			for k, v := range rd.rows {
				fmt.Fprintf(&buf, "%s:  %s,\n", k, string((v)[:]))
			}
			io.WriteString(s, buf.String())
			return
		}

		fallthrough
	case 's', 'q':
		var buf bytes.Buffer
		for k, v := range rd.rows {
			fmt.Fprintf(&buf, "%s:%s,", k, string((v)[:]))
		}
		io.WriteString(s, buf.String())
	}
}

func (rd *ActiveRecord) Set(colName string, b []byte) error {
	if colName == "" {
		return errors.New("col name is empty")
	}
	rd.rows[colName] = b
	return nil
}

func (rd *ActiveRecord) SetString(colName string, v string) error {
	return rd.Set(colName, []byte(v))
}

func (rd *ActiveRecord) SetInt(colName string, v int) error {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	return rd.Set(colName, b)
}

func (rd *ActiveRecord) SetInt64(colName string, v int64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return rd.Set(colName, b)
}
