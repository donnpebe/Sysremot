package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func swapJob(roundedTime time.Time) {
	conn := pool.Get()
	defer conn.Close()

	currentKey := fmt.Sprintf("%s|swap|current", AppName)
	historyKey := fmt.Sprintf("%s|swap|%s|%s", AppName, roundedTime.Format("2006-01-02"), roundedTime.Format("15:04:05"))

	swap := sigar.Swap{}
	err := swap.Get()
	if err != nil {
		errLogger.Println("SwapJob error: ", err)
		return
	}

	data := fmt.Sprintf(`{"total":"%d","used":"%d","free":"%d","in-percent":"%d"}`,
		swap.Total, swap.Used, swap.Free, usePercent(swap.Total, swap.Free))

	conn.Send("MULTI")
	conn.Send("SETEX", currentKey, TheTicker.Seconds()*2, data)
	conn.Send("SETEX", historyKey, ExpireInterval, data)
	_, err = conn.Do("EXEC")
	if err != nil {
		errLogger.Println(err)
	}
}
