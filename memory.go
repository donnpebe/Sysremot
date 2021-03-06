package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func memoryJob(roundedTime time.Time) {
	conn := pool.Get()
	defer conn.Close()

	currentKey := fmt.Sprintf("%s|memory|current", AppName)
	historyKey := fmt.Sprintf("%s|memory|%s|%s", AppName, roundedTime.Format("2006-01-02"), roundedTime.Format("15:04:05"))

	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		errLogger.Println("MemoryJob error: ", err)
		return
	}

	data := fmt.Sprintf(`{"total":"%d","used":"%d","free":"%d","in-percent":"%d"}`,
		mem.Total, mem.ActualUsed, mem.ActualFree, usePercent(mem.Total, mem.ActualFree))

	conn.Send("MULTI")
	conn.Send("SETEX", currentKey, TheTicker.Seconds()*2, data)
	conn.Send("SETEX", historyKey, ExpireInterval, data)
	_, err = conn.Do("EXEC")
	if err != nil {
		errLogger.Println(err)
	}
}
