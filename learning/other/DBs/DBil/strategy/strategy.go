package strategy

type Strategy interface {
	Save(key, value string) error
	Get(key string) (value *string)
	Index() error
	Clean() error
}
