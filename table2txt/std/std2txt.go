package std2txt

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	t2t "github.com/takanoriyanagitani/go-pg2txt/table2txt"
)

const TimeoutDefault time.Duration = 10 * time.Second

const SchemaDefault string = "public"

type RowsToText interface {
	Convert(context.Context, *sql.Rows) ([]string, error)
}

type Rows2txtFn func(context.Context, *sql.Rows) ([]string, error)

func (f Rows2txtFn) Convert(c context.Context, r *sql.Rows) ([]string, error) {
	return f(c, r)
}

func (f Rows2txtFn) AsIf() RowsToText { return f }

func TryForEach(r *sql.Rows, cb func(*sql.Rows) error) error {
	for r.Next() {
		e := cb(r)
		if nil != e {
			return e
		}
	}
	return nil
}

func Rows2txtFnFromTimeout(timeout time.Duration) Rows2txtFn {
	return func(ctx context.Context, r *sql.Rows) (rows []string, e error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var buf string
		e = TryForEach(
			r,
			func(ro *sql.Rows) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				e := ro.Scan(&buf)
				if nil != e {
					return e
				}
				rows = append(rows, buf)
				return nil
			},
		)
		return rows, e
	}
}

var Rows2txtFnDefault Rows2txtFn = Rows2txtFnFromTimeout(TimeoutDefault)

type Tab2Jsons struct {
	db  *sql.DB
	r2t RowsToText
}

func (j Tab2Jsons) Query(c context.Context, t string) ([]string, error) {
	var query string = fmt.Sprintf(
		`
            SELECT
                ROW_TO_JSON(t) AS json_string
            FROM %s AS t
        `,
		t,
	)
	rows, e := j.db.QueryContext(c, query)
	if nil != e {
		return nil, e
	}
	defer rows.Close()

	return j.r2t.Convert(c, rows)
}

func (j Tab2Jsons) AsIf() t2t.TableToString { return j }

func (j Tab2Jsons) WithRowsToText(r2t RowsToText) Tab2Jsons {
	j.r2t = r2t
	return j
}

func (j Tab2Jsons) WithTimeout(timeout time.Duration) Tab2Jsons {
	var r2t RowsToText = Rows2txtFnFromTimeout(timeout).AsIf()
	return j.WithRowsToText(r2t)
}

func Tab2JsonsNew(db *sql.DB) Tab2Jsons {
	return Tab2Jsons{
		db:  db,
		r2t: Rows2txtFnDefault.AsIf(),
	}
}

type TabChkInfo struct {
	db     *sql.DB
	schema string
}

func (i TabChkInfo) Check(c context.Context, tableName string) error {
	const query string = `
        SELECT table_name
        FROM information_schema.tables
        WHERE
          table_schema = $1
          AND table_name = $2
    `
	var row *sql.Row = i.db.QueryRowContext(
		c,
		query,
		i.schema,
		tableName,
	)
	var buf string
	e := row.Scan(&buf)
	if nil != e {
		return e
	}
	if buf != tableName {
		return fmt.Errorf("unexpected table name: %s", buf)
	}
	return nil
}

func (i TabChkInfo) WithSchema(schema string) TabChkInfo {
	i.schema = schema
	return i
}

func (i TabChkInfo) AsIf() t2t.TableChecker { return i }

func TabChkInfoNew(db *sql.DB) TabChkInfo {
	return TabChkInfo{
		db:     db,
		schema: SchemaDefault,
	}
}
