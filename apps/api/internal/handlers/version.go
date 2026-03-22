package handlers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"pocketpanel/api/internal/models"
)

type VersionsHandler struct {
	db *gorm.DB
}

func NewVersionsHandler(db *gorm.DB) *VersionsHandler {
	return &VersionsHandler{
		db: db,
	}
}

func (h *VersionsHandler) GetVersions(c fiber.Ctx) error {
	serverType := models.ServerType(c.Params("type"))

	if serverType != models.ServerTypeVanilla && serverType != models.ServerTypeFabric {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server type. Must be 'vanilla' or 'fabric'",
		})
	}

	var versions []models.Version
	if err := h.db.Where("server_type = ?", serverType).Order("version DESC").Find(&versions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch versions",
		})
	}

	if len(versions) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No versions found. Try syncing versions first.",
		})
	}

	versionStrings := make([]string, len(versions))
	for i, v := range versions {
		versionStrings[i] = v.Version
	}

	sort.Slice(versionStrings, func(i, j int) bool {
		return compareVersions(versionStrings[j], versionStrings[i]) < 0
	})

	return c.JSON(fiber.Map{
		"server_type": serverType,
		"versions":    versionStrings,
	})
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, _ := strconv.Atoi(aParts[i])
		bNum, _ := strconv.Atoi(bParts[i])
		if aNum != bNum {
			if aNum > bNum {
				return 1
			}
			return -1
		}
	}
	return len(aParts) - len(bParts)
}
