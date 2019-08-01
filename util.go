package trustar

import "time"

// TimeToMsEpoch converts time.Time to milliseconds epoch string
func TimeToMsEpoch(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// MsEpochToTime converts millisecond epoch int to time.Time
func MsEpochToTime(msInt int64) (time.Time, error) {

	millisPerSecond := int64(time.Second / time.Millisecond)
	nanosPerMillisecond := int64(time.Millisecond / time.Nanosecond)

	return time.Unix(msInt/millisPerSecond,
		(msInt%millisPerSecond)*nanosPerMillisecond), nil
}
