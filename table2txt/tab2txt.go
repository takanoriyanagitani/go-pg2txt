package tab2txt

import (
	"context"
)

type TableToString interface {
	Query(ctx context.Context, tableName string) (rows []string, e error)
}

type Tab2StrFn func(context.Context, string) ([]string, error)

func (f Tab2StrFn) Query(c context.Context, t string) ([]string, error) {
	return f(c, t)
}
func (f Tab2StrFn) AsIf() TableToString { return f }

type TableChecker interface {
	Check(ctx context.Context, tableName string) error
}

type TabChkFn func(context.Context, string) error

func (f TabChkFn) Check(c context.Context, t string) error { return f(c, t) }
func (f TabChkFn) AsIf() TableChecker                      { return f }

func (f TabChkFn) ToChecked(t2s TableToString) TableToString {
	return Tab2StrFn(func(c context.Context, t string) ([]string, error) {
		e := f.Check(c, t)
		if nil != e {
			return nil, e
		}
		return t2s.Query(c, t)
	}).AsIf()
}
