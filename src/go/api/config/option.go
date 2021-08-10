package config

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeJSON
	DataTypeYAML
)

type CreateOption func(*createOptions)

type createOptions struct {
	path     string
	data     []byte
	dataType DataType
	validate bool
}

func newCreateOptions(opts ...CreateOption) createOptions {
	var o createOptions

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func CreateFromPath(p string) CreateOption {
	return func(o *createOptions) {
		o.path = p
	}
}

func CreateFromJSON(d []byte) CreateOption {
	return func(o *createOptions) {
		o.data = d
		o.dataType = DataTypeJSON
	}
}

func CreateFromYAML(d []byte) CreateOption {
	return func(o *createOptions) {
		o.data = d
		o.dataType = DataTypeYAML
	}
}

func CreateWithValidation() CreateOption {
	return func(o *createOptions) {
		o.validate = true
	}
}
