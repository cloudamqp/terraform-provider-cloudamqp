package durations

import "time"

func FiveSeconds() time.Duration {
	return 5 * time.Second
}

func TenSeconds() time.Duration {
	return 10 * time.Second
}

func ThirtySeconds() time.Duration {
	return 30 * time.Second
}

func OneMinute() time.Duration {
	return time.Minute
}

func FiveMinutes() time.Duration {
	return 5 * time.Minute
}

func TenMinutes() time.Duration {
	return 10 * time.Minute
}

func Hour() time.Duration {
	return time.Hour
}
