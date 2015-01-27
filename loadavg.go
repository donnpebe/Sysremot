package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func loadAvgJob(roundedTime time.Time) {
	conn := pool.Get()
	defer conn.Close()

	currentKey := fmt.Sprintf("%s|loadavg|current", AppName)

	load := sigar.LoadAverage{}
	err := load.Get()
	if err != nil {
		errLogger.Println("LoadAvgJob error: ", err)
		return
	}

	data := fmt.Sprintf(`{"one":"%.2f","five":"%.2f","fifteen":"%.2f"}`,
		load.One, load.Five, load.Fifteen)

	conn.Send("MULTI")
	conn.Send("SETEX", currentKey, TheTicker.Seconds()*2, data)
	_, err = conn.Do("EXEC")
	if err != nil {
		errLogger.Println(err)
	}
}
