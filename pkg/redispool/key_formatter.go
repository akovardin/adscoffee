package redispool

type KeyFormatter interface {
	FormatKey(key string) string
}

type DummyKeyFormatter struct{}

func (DummyKeyFormatter) FormatKey(key string) string {
	return key
}
