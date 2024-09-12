package utils

import "time"

func AdjustTimeToUserTimezone(t time.Time, timezone string) (time.Time, error) {
	t = t.UTC()
	userLocation, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	adjustedTime := t.In(userLocation)
	return adjustedTime, nil
}
