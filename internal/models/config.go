package models

type Config interface {
	InitEnv()
	ReadEnv()
	Get(string) string
}
