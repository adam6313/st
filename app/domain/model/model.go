package model

// Modes -
type Mode int32

const (
	// 開發模式
	Mode_Development Mode = iota

	// 正式模式
	Mode_Production
)

var Mode_name = map[int32]string{
	0: "development",
	1: "production",
}

var Mode_value = map[string]int32{
	"development": 0,
	"production":  1,
}

// ModelVerify - 驗證模式
func (m *Mode) ModelVerify(s string) Mode {
	v, ok := Mode_value[s]
	if ok {
		return Mode(v)
	}

	panic("Start mode not allowed")
}
