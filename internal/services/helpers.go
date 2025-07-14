package services

import (
	"github.com/Lekuruu/go-puush/internal/app"
	"gorm.io/gorm"
)

func preloadQuery(state *app.State, preload []string) *gorm.DB {
	result := state.Database

	for _, p := range preload {
		result = result.Preload(p)
	}

	return result
}
