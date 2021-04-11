package prot

type protOptions struct {
	rowSep    byte
	columnSep byte
}

type funcProtOption struct {
	f func(*protOptions)
}
type ProtOption interface {
	apply(*protOptions)
}

func (fpo *funcProtOption) apply(o *protOptions) {
	fpo.f(o)
}

func defaultOpt() protOptions {
	return protOptions{
		rowSep:    '\n',
		columnSep: ',',
	}
}

// WithRowSep 修改行分隔符
func WithRowSep(sep byte) *funcProtOption {
	return &funcProtOption{func(opt *protOptions) {
		opt.rowSep = sep
	}}
}

// WithColumnSep 修改列分隔符
func WithColumnSep(sep byte) *funcProtOption {
	return &funcProtOption{func(opt *protOptions) {
		opt.columnSep = sep
	}}
}
