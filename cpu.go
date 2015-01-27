package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func cpuJob(start time.Time) {
	conn := pool.Get()
	defer conn.Close()

	// round the time to the closest minute (round down)
	roundedTs := roundTheTimestamp(start.Unix(), int64(TheTicker.Seconds()))
	// convert back to time.Time
	roundedTime := time.Unix(roundedTs, 0)
	// round the time to the closest hour (round down) then add the expire time
	expireTs := roundTheTimestamp(start.Unix(), ExpireInterval/2) + ExpireInterval
	cpuCountKey := fmt.Sprintf("%s|cpu|count", AppName)

	cpulist := sigar.CpuList{}
	err := cpulist.Get()
	if err != nil {
		errLogger.Println("CpuJob error: ", err)
		return
	}

	_, err = conn.Do("SET", cpuCountKey, len(cpulist.List))
	if err != nil {
		errLogger.Println("CpuJob error: ", err)
		return
	}

	for i, cpu := range cpulist.List {
		total := cpu.Total()
		used := total - cpu.Idle
		data := fmt.Sprintf(`{"total":"%d","used":"%d","free":"%d","in-percent":"%d"}`,
			total, used, cpu.Idle, usePercent(total, cpu.Idle))

		currentKey := fmt.Sprintf("%s|cpu:%d|current", AppName, i)
		historyKey := fmt.Sprintf("%s|cpu:%d|%s", AppName, i, roundedTime.Format("2006-01-02"))
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
}
