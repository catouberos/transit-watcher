package handler

import (
	"errors"
	"strings"
	"time"
)

const (
	OpTimeLayout = "15:04"
)

func ParseOperationTime(operationTime string) (from, to time.Time, err error) {
	unparsedTimes := strings.Split(operationTime, " - ")

	if len(unparsedTimes) != 2 {
		// nothing to parse
		return time.Time{}, time.Time{}, errors.New("Operation time is not supported")
	}

	now := time.Now()

	from, err = time.Parse(OpTimeLayout, unparsedTimes[0])

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	from = time.Date(now.Year(), now.Month(), now.Day(), from.Hour(), from.Minute(), from.Second(), 0, now.Location())

	to, err = time.Parse(OpTimeLayout, unparsedTimes[1])

	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	to = time.Date(now.Year(), now.Month(), now.Day(), from.Hour(), from.Minute(), from.Second(), 0, now.Location())

	return from, to, nil
}
