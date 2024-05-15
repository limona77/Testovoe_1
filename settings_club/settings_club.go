package settings_club

import "time"

type SettingsCLub struct {
	CountTables  int
	StartTime    time.Time
	EndTime      time.Time
	PricePerHour int
}

func NewSettingsClub(countTables int, startTime, endTime time.Time, pricePerHour int) SettingsCLub {
	return SettingsCLub{
		CountTables:  countTables,
		StartTime:    startTime,
		EndTime:      endTime,
		PricePerHour: pricePerHour,
	}
}
