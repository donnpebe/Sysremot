package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func uptimeJob(roundedTime time.Time) {
	conn := pool.Get()
	defer conn.Close()

	currentKey := fmt.Sprintf("%s|uptime|current", AppName)

	uptime := sigar.Uptime{}
	err := uptime.Get()
	if err != nil {
		errLogger.Println("UptimeJob error: ", err)
		return
	}

	data := fmt.Sprintf(`{"length":"%d"}`, uint64(uptime.Length))

	conn.Send("MULTI")
	conn.Send("SETEX", currentKey, TheTicker.Seconds()*2, data)
	_, err = conn.Do("EXEC")
	if err != nil {
		errLogger.Println(err)
	}
}
