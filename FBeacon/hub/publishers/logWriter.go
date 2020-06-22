package publishers

import (
	"fmt"

	"github.com/scottmccool/FBeacon/readings"
)

func PublishLog(c chan readings.FBeacon) {
	for {
		select {
		case obs := <-c:
			fmt.Println(obs)
		}
	}
}
