package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func swapJob(start time.Time) {
	conn := pool.Get()
	defer conn.Close()

	// round the time to the closest minute (round down)
	roundedTs := roundTheTimestamp(start.Unix(), int64(TheTicker.Seconds()))
	// convert back to time.Time
	roundedTime := time.Unix(roundedTs, 0)
	// round the time to the closest hour (round down) then add the expire time
	expireTs := roundTheTimestamp(start.Unix(), ExpireInterval/2) + ExpireInterval
	currentKey := fmt.Sprintf("%s|swap|current", AppName)
	historyKey := fmt.Sprintf("%s|swap|%s", AppName, roundedTime.Format("2006-01-02"))

	swap := sigar.Swap{}
	err := swap.Get()
	if err != nil {
		errLogger.Println("SwapJob error: ", err)
		return
	}

	data := fmt.Sprintf(`{"total":"%d","used":"%d","free":"%d","in-percent":"%d"}`,
		swap.Total, swap.Used, swap.Free, usePercent(swap.Total, swap.Free))

	field := roundedTime.Format("15:04:05")

	conn.Send("MULTI")
	conn.Send("SETEX", currentKey, TheTicker.Seconds()*2, data)
	conn.Send("HSET", historyKey, field, data)
	conn.Send("EXPIREAT", historyKey, expireTs)
	_, err = conn.Do("EXEC")
	if err != nil {
		errLogger.Println(err)
	}
}
