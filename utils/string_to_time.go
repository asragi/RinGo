package utils

import "time"

func StringToTime(utcTime string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"
	t, err := time.Parse(layout, utcTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
