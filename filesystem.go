package main

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
)

func formatSize(size uint64) uint64 {
	return size * 1024
}

func filesystemJob(start time.Time) {
	conn := pool.Get()
	defer conn.Close()

	// round the time to the closest minute (round down)
	roundedTs := roundTheTimestamp(start.Unix(), int64(TheTicker.Seconds()))
	// convert back to time.Time
	roundedTime := time.Unix(roundedTs, 0)
	// round the time to the closest hour (round down) then add the expire time
	expireTs := roundTheTimestamp(start.Unix(), ExpireInterval/2) + ExpireInterval

	fslist := sigar.FileSystemList{}
	err := fslist.Get()
	if err != nil {
		errLogger.Println("FilesystemJob error: ", err)
		return
	}

	for _, fs := range fslist.List {
		dirname := fs.DirName
		if dirname != "/" && dirname != "/home" {
			continue
		}

		usage := sigar.FileSystemUsage{}
		usage.Get(dirname)

		if usage.Total == 0 {
			continue
		}

		data := fmt.Sprintf(`{"total":"%d","used":"%d","free":"%d","in-percent":"%d"}`,
			formatSize(usage.Total), formatSize(usage.Used), formatSize(usage.Avail), usePercent(formatSize(usage.Total), formatSize(usage.Avail)))

		currentKey := fmt.Sprintf("%s|fs:%s|current", AppName, dirname)
		historyKey := fmt.Sprintf("%s|fs:%s|%s", AppName, dirname, roundedTime.Format("2006-01-02"))
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
