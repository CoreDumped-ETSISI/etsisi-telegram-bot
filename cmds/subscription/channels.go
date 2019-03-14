package subscription

var (
	publicChannels = []string{
		"noticias",
		"eventos",
		"avisos",
		"status",
		"coredumped",
		"menu",
	}

	normalMap = map[string]string{
		"noticias":   "UNIVERSITY_NEWS_CHANNEL",
		"eventos":    "UNIVERSITY_EVENTOS_CHANNEL",
		"avisos":     "UNIVERSITY_AVISOS_CHANNEL",
		"status":     "SERVICE_STATUS_CHANNEL",
		"status_bot": "SERVICE_STATUS_CHANNEL_INTERNAL",
		"coredumped": "UNIVERSITY_COREDUMPED_CHANNEL",
		"menu":       "CAFETERIA_NEW_MENU_AVAILABLE",
	}
	reverseMap = map[string]string{
		"UNIVERSITY_NEWS_CHANNEL":         "noticias",
		"UNIVERSITY_EVENTOS_CHANNEL":      "eventos",
		"UNIVERSITY_AVISOS_CHANNEL":       "avisos",
		"SERVICE_STATUS_CHANNEL":          "status",
		"SERVICE_STATUS_CHANNEL_INTERNAL": "status_bot",
		"UNIVERSITY_COREDUMPED_CHANNEL":   "coredumped",
		"CAFETERIA_NEW_MENU_AVAILABLE":    "menu",
	}

	redisChannels = []string{
		"UNIVERSITY_NEWS_CHANNEL",
		"UNIVERSITY_EVENTOS_CHANNEL",
		"UNIVERSITY_AVISOS_CHANNEL",
		"SERVICE_STATUS_CHANNEL",
		"SERVICE_STATUS_CHANNEL_INTERNAL",
		"UNIVERSITY_COREDUMPED_CHANNEL",
		"CAFETERIA_NEW_MENU_AVAILABLE",
	}
)
