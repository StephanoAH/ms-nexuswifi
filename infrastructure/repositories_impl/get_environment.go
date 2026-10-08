package repositories_impl

import (
	"fmt"
	"os"
	"strconv"
)

type GetEnvironmentImpl struct {
}

func (g *GetEnvironmentImpl) GetEnvironmentString(key string) string {
	result := os.Getenv(key)
	return result
}

func (g *GetEnvironmentImpl) GetEnvironmentInt(key string) (int, error) {
	result, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0, fmt.Errorf("error convirtiendo %s a entero: %v", key, err)
	}
	return result, nil
}

func (g *GetEnvironmentImpl) GetEnvironmentBool(key string) (bool, error) {
	result, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false, fmt.Errorf("error convirtiendo %s a booleano: %v", key, err)
	}
	return result, nil
}
