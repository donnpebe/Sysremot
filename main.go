package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/takama/daemon"
)

const (
	// AppName namespace in redis key
	AppName = "sysremot"
	// TheTicker theticker control how often worker do the resource gathering
	TheTicker = 1 * time.Minute
	// ExpireInterval value that determine who long history data be keep in redis
	ExpireInterval = 7200
	// Description describe the app
	Description    = "System Resource Monitoring Tool"
	rootPrivileges = "You must have root user privileges. Possibly using 'sudo' command should help"
)

// Job define type of job to do in worker
type Job func(start time.Time)

var (
	pool       *redis.Pool
	errLogger  = log.New(os.Stderr, "", log.LstdFlags)
	stdLogger  = log.New(os.Stdout, "", log.LstdFlags)
	configdir  = fmt.Sprintf("/etc/%s", AppName)
	configfile = fmt.Sprintf("%s/%s.env", configdir, AppName)
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

func (service *Service) createConfigFile() (bool, error) {
	if checkPrivileges() == false {
		return false, errors.New(rootPrivileges)
	}

	_, err := os.Stat(configdir)
	if err != nil {
		err = os.Mkdir(configdir, 0644)
		if err != nil {
			return false, err
		}
	}

	_, err = os.Stat(configfile)
	if err == nil {
		stdLogger.Println("Config file already exist")
		return true, nil
	}

	file, err := os.Create(configfile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Write(configTemplate)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (service *Service) removeConfigFile() (bool, error) {
	if checkPrivileges() == false {
		return false, errors.New(rootPrivileges)
	}

	_, err := os.Stat(configfile)
	if err == nil {
		if err = os.Remove(configfile); err != nil {
			return false, err
		}
		if err = os.Remove(configdir); err != nil {
			return false, err
		}
		return true, nil
	}

	return true, nil

}

// Manage entrypoint to managing app
func (service *Service) Manage() (string, error) {
	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", AppName)

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			_, err := service.createConfigFile()
			if err != nil {
				return "", err
			}
			return service.Install()
		case "remove":
			_, err := service.removeConfigFile()
			if err != nil {
				return "", err
			}
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Get environment variable in .env
	err := godotenv.Load(fmt.Sprintf("/etc/%s/%s.env", AppName, AppName))
	if err != nil {
		errLogger.Fatalln("Error loading .env file", err)
	}

	poolSizeStr := os.Getenv("SRMT_REDIS_POOL_SIZE")
	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		errLogger.Fatalln("wrong pool size value")
	}

	// initialize redis pool to be used in worker
	pool = redisPool(os.Getenv("SRMT_REDIS_SERVER"), os.Getenv("SRMT_REDIS_PASS"), poolSize)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// make the job run forever until killed manually
	go func() {
		for t := range time.Tick(TheTicker) {
			stdLogger.Println("Cooking metrics started")
			jobs := make([]Job, 6)
			jobs[0] = uptimeJob
			jobs[1] = memoryJob
			jobs[2] = swapJob
			jobs[3] = cpuJob
			jobs[4] = loadAvgJob
			jobs[5] = filesystemJob

			for _, job := range jobs {
				go job(t)
			}
		}
	}()
	stdLogger.Printf("%s is ready to cook...\n", AppName)

	killSignal := <-interrupt
	stdLogger.Println("Got signal:", killSignal)
	if killSignal == os.Interrupt {
		stdLogger.Printf("%s was interruped by system signal\n", AppName)
	}

	return fmt.Sprintf("R.I.P %s", AppName), nil
}

func main() {
	// initialize daemon for this app
	srv, err := daemon.New(AppName, Description)
	if err != nil {
		errLogger.Fatalf("Error initializing daemon: \n", err)
	}
	service := &Service{srv}

	status, err := service.Manage()
	if err != nil {
		errLogger.Fatalln(err)
	}
	stdLogger.Println(status)
}