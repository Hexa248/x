package controllers

import (
	"sort"

	"hotel-booking/models"
)

func notificationsByRole(all []models.Notification, role models.Role) []models.Notification {
	filtered := make([]models.Notification, 0)
	for _, n := range all {
		if n.Role == role {
			filtered = append(filtered, n)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})
	if len(filtered) > 8 {
		return filtered[:8]
	}
	return filtered
}
