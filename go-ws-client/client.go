package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Report struct {
	Time                     int
	ConnectionDurations      map[string]int
	PingPoingEventDurations  map[string]int
	TotalConnections         int
	DroppedConnections       int
	ReceivedPingPongMessages int
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
func initializeDurationMaps(durationsMap *map[string]int) {
	(*durationsMap) = make(map[string]int)
	(*durationsMap)["under_00050"] = 0
	(*durationsMap)["under_00100"] = 0
	(*durationsMap)["under_00200"] = 0
	(*durationsMap)["under_00500"] = 0
	(*durationsMap)["under_00750"] = 0
	(*durationsMap)["under_01000"] = 0
	(*durationsMap)["under_01250"] = 0
	(*durationsMap)["under_01500"] = 0
	(*durationsMap)["under_01750"] = 0
	(*durationsMap)["under_02000"] = 0
	(*durationsMap)["under_03000"] = 0
	(*durationsMap)["under_04000"] = 0
	(*durationsMap)["under_05000"] = 0
	(*durationsMap)["under_06000"] = 0
	(*durationsMap)["under_07000"] = 0
	(*durationsMap)["under_08000"] = 0
	(*durationsMap)["under_09000"] = 0
	(*durationsMap)["under_10000"] = 0
	(*durationsMap)["zzzzz_10000"] = 0
}
func analyzeDuration(durationsMap *map[string]int, duration int64) {
	if duration < 50 {
		(*durationsMap)["under_00050"]++
	} else if duration < 100 {
		(*durationsMap)["under_00100"]++
	} else if duration < 200 {
		(*durationsMap)["under_00200"]++
	} else if duration < 500 {
		(*durationsMap)["under_00500"]++
	} else if duration < 750 {
		(*durationsMap)["under_00750"]++
	} else if duration < 1000 {
		(*durationsMap)["under_01000"]++
	} else if duration < 1250 {
		(*durationsMap)["under_01250"]++
	} else if duration < 1500 {
		(*durationsMap)["under_01500"]++
	} else if duration < 1750 {
		(*durationsMap)["under_01750"]++
	} else if duration < 2000 {
		(*durationsMap)["under_02000"]++
	} else if duration < 3000 {
		(*durationsMap)["under_03000"]++
	} else if duration < 4000 {
		(*durationsMap)["under_04000"]++
	} else if duration < 5000 {
		(*durationsMap)["under_05000"]++
	} else if duration < 6000 {
		(*durationsMap)["under_06000"]++
	} else if duration < 7000 {
		(*durationsMap)["under_07000"]++
	} else if duration < 8000 {
		(*durationsMap)["under_08000"]++
	} else if duration < 9000 {
		(*durationsMap)["under_09000"]++
	} else if duration < 10000 {
		(*durationsMap)["under_10000"]++
	} else {
		(*durationsMap)["zzzzz_10000"]++
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var mutex = &sync.Mutex{}
	var report = &Report{}
	initializeDurationMaps(&report.ConnectionDurations)
	initializeDurationMaps(&report.PingPoingEventDurations)
	var showEstablishedConnections = true
	var getCachedEvents = false
	var currentTotalConnections = 0
	var portIndex = 0
	totalConnections, _ := strconv.Atoi(os.Args[1])
	batchSize, _ := strconv.Atoi(os.Args[2])
	mode, _ := strconv.Atoi(os.Args[3])

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// mode 0 => LB
	// mode 1 => VM
	go func() {
		for showEstablishedConnections {
			// log.Println("ESTABLISHED CONNECTIONS: " + strconv.Itoa(report.TotalConnections))
			time.Sleep(2 * time.Second)
		}
	}()
	timing := time.Now()
	for totalConnections > currentTotalConnections {
		currentTotalConnections = currentTotalConnections + batchSize
		for i := 1; i <= batchSize; i++ {
			go func() {
				connectionStart := time.Now()
				getCacheEventStart := time.Now()
				port := 10000
				mutex.Lock()
				port = port + (portIndex % 10)
				portIndex++
				mutex.Unlock()
				host := "localhost"
				isSsl := false
				url := ""
				if mode < 1 {
					port = 5000
					url = gosocketio.GetUrl(host, port, isSsl)
				} else {
					url = "ws://localhost:30080/wsk/socket.io/?EIO=3&transport=websocket"
					isSsl = false
				}
				c, err := gosocketio.Dial(
					url,
					&transport.WebsocketTransport{
						PingInterval:   3 * time.Second,
						PingTimeout:    15 * time.Second,
						ReceiveTimeout: 15 * time.Second,
						SendTimeout:    15 * time.Second,
						BufferSize:     1024 * 32,
					})

				if err == nil {
					err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
						mutex.Lock()
						report.TotalConnections--
						report.DroppedConnections++
						mutex.Unlock()
					})
					err = c.On("pong", func(h *gosocketio.Channel) {
						log.Println("pong")
					})
					err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
						connectionDuration := time.Since(connectionStart)
						mutex.Lock()
						report.TotalConnections++
						analyzeDuration(&report.ConnectionDurations, connectionDuration.Milliseconds())
						mutex.Unlock()
						go func() {
							var getCachedEventsLocal = true
							for {
								if getCachedEvents && getCachedEventsLocal {
									getCachedEventsLocal = false
									getCacheEventStart = time.Now()
									if err := c.Emit("test000", "me"); err != nil {
										break
									}
								}
								time.Sleep(5 * time.Second)
							}
						}()
					})
					err = c.On("okok", func(h *gosocketio.Channel) {
						cacheEventDuration := time.Since(getCacheEventStart)
						mutex.Lock()
						report.ReceivedPingPongMessages++
						analyzeDuration(&report.PingPoingEventDurations, cacheEventDuration.Milliseconds())
						mutex.Unlock()
					})
				}
			}()
		}
		time.Sleep(500 * time.Millisecond)
	}
	showEstablishedConnections = false
	// log.Println("DONE WITH ESTABLISHING CONNECTIONS: " + strconv.Itoa(report.TotalConnections))
	// log.Println("Press enter to send ping pong messages...")
	go func() {
		for {
			report.Time = int(time.Since(timing).Seconds())
			log.Println(prettyPrint(report))
			time.Sleep(2 * time.Second)
		}
	}()
	var dummy string
	fmt.Scanln(&dummy)
	getCachedEvents = true
	log.Println("Press enter to stop...")
	fmt.Scanln(&dummy)
}
