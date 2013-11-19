package orm

import (
	"fmt"
	"strconv"
)

var oracleOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "= UPPER(?)",
	"contains":    "LIKE ?",
	"icontains":   "LIKE UPPER(?)",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"startswith":  "LIKE ?",
	"endswith":    "LIKE ?",
	"istartswith": "LIKE UPPER(?)",
	"iendswith":   "LIKE UPPER(?)",
}

var oracleTypes = map[string]string{
	"auto":            "serial NOT NULL PRIMARY KEY",
	"pk":              "NOT NULL PRIMARY KEY",
	"bool":            "bool",
	"string":          "varchar(%d)",
	"string-text":     "text",
	"time.Time-date":  "date",
	"time.Time":       "timestamp with time zone",
	"int8":            `smallint CHECK("%COL%" >= -127 AND "%COL%" <= 128)`,
	"int16":           "smallint",
	"int32":           "integer",
	"int64":           "bigint",
	"uint8":           `smallint CHECK("%COL%" >= 0 AND "%COL%" <= 255)`,
	"uint16":          `integer CHECK("%COL%" >= 0)`,
	"uint32":          `bigint CHECK("%COL%" >= 0)`,
	"uint64":          `bigint CHECK("%COL%" >= 0)`,
	"float64":         "double precision",
	"float64-decimal": "numeric(%d, %d)",
}

type dbBaseOracle struct {
	dbBase
}

var _ dbBaser = new(dbBaseOracle)

func (d *dbBaseOracle) OperatorSql(operator string) string {
	return oracleOperators[operator]
}

func (d *dbBaseOracle) GenerateOperatorLeftCol(fi *fieldInfo, operator string, leftCol *string) {
	switch operator {
	case "contains", "startswith", "endswith":
		*leftCol = fmt.Sprintf("%s::text", *leftCol)
	case "iexact", "icontains", "istartswith", "iendswith":
		*leftCol = fmt.Sprintf("UPPER(%s::text)", *leftCol)
	}
}

func (d *dbBaseOracle) SupportUpdateJoin() bool {
	return false
}

func (d *dbBaseOracle) MaxLimit() uint64 {
	return 0
}

func (d *dbBaseOracle) TableQuote() string {
	return `"`
}

func (d *dbBaseOracle) ReplaceMarks(query *string) {
	q := *query
	num := 0
	for _, c := range q {
		if c == '?' {
			num += 1
		}
	}
	if num == 0 {
		return
	}
	data := make([]byte, 0, len(q)+num)
	num = 1
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '?' {
			data = append(data, ':')
			data = append(data, []byte(strconv.Itoa(num))...)
			num += 1
		} else {
			data = append(data, c)
		}
	}
	*query = string(data)
}

func (d *dbBaseOracle) HasReturningID(mi *modelInfo, query *string) (has bool) {
	if mi.fields.pk.auto {
		if query != nil {
			*query = fmt.Sprintf(`%s RETURNING "%s"`, *query, mi.fields.pk.column)
		}
		has = true
	}
	return
}

func (d *dbBaseOracle) ShowTablesQuery() string {
	return "select table_name from user_tables"
}

func (d *dbBaseOracle) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM user_tab_columns  where  table_name = '%s'", table)
}

func (d *dbBaseOracle) DbTypes() map[string]string {
	return oracleTypes
}

func (d *dbBaseOracle) IndexExists(db dbQuerier, table string, name string) bool {
	query := fmt.Sprintf("select count(t.*),i.index_type from user_ind_columns t,user_indexes i where t.index_name = i.%s and t.table_name='%s'", name, table)
	row := db.QueryRow(query)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

func newdbBaseOracle() dbBaser {
	b := new(dbBaseOracle)
	b.ins = b
	return b
}
