package config

import (
	"os"
	"strconv"
)

type Config struct {
	Addr           string
	PythonURL      string
	DataDir        string
	APIVersion     string
}

func Load() Config {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	py := os.Getenv("PYTHON_TRANSFORM_URL")
	if py == "" {
		py = "http://127.0.0.1:5000"
	}
	data := os.Getenv("DATA_DIR")
	if data == "" {
		data = "data"
	}
	return Config{
		Addr:       addr,
		PythonURL:  py,
		DataDir:    data,
		APIVersion: "v1",
	}
}

func BoolEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
