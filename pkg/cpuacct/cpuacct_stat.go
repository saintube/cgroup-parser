/*
Copyright © 2022 saintube

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cpuacct

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/saintube/cgroup-parser/pkg/utils"
)

const StatFileName = "cpuacct.stat"

var StatKeys = map[string]bool{"user": true, "system": true}

func Stat(cgroupPath string, interval int, duration int) error {
	// workflow
	// 1. validate arguments
	// 2. get system Jifies
	// 3. loop reading cpuacct.stat and calculating the usage by intervals
	log.Println("start parse cpuacct.stat")

	err := checkStatArgs(cgroupPath, interval, duration)
	if err != nil {
		return err
	}

	j, err := getJiffies()
	if err != nil {
		return err
	}

	return loopReadStat(cgroupPath, interval, duration, j)
}

func checkStatArgs(cgroupPath string, interval int, duration int) error {
	err := utils.FilePathExists(filepath.Join(cgroupPath, StatFileName))
	if err != nil {
		return utils.WrapError(err.Error(), "invalid cgroup-path")
	}
	if interval <= 0 {
		return fmt.Errorf("interval should be larger than zero")
	}
	if duration <= 0 {
		return fmt.Errorf("duration should be larger than zero")
	}
	return nil
}

func getJiffies() (float64, error) {
	out, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		return -1, utils.WrapError(err.Error(), "failed to get jiffies")
	}
	v, err := strconv.ParseInt(strings.Trim(string(out), "\n"), 10, 64)
	if err != nil {
		return -1, utils.WrapError(err.Error(), "failed to get jiffies")
	}
	if v <= 0 {
		return -1, fmt.Errorf("invalid jiffies %v", v)
	}
	j := 1.0 / float64(v)
	log.Printf("get jiffies %v (seconds)", j)
	return j, nil
}

func loopReadStat(cgroupPath string, interval, duration int, jiffies float64) error {
	var count int64
	var ts time.Time
	nrReads := duration / interval
	for i := 0; i < nrReads; i++ {
		fd, err := ioutil.ReadFile(filepath.Join(cgroupPath, StatFileName))
		curTs := time.Now()
		if err != nil {
			return utils.WrapError(err.Error(), "failed to read cgroup file")
		}
		lines := strings.Split(strings.TrimSpace(string(fd)), "\n")
		if len(lines) <= 0 || len(lines) > 2 {
			return fmt.Errorf("invalid rows for cpuacct.stat content: %s", string(fd))
		}
		curCount := int64(0)
		for _, line := range lines {
			ss := strings.Fields(line)
			if len(ss) != 2 || !StatKeys[ss[0]] {
				return fmt.Errorf("invalid columns for cpuacct.stat content: %v", lines)
			}
			ticks, err := strconv.ParseInt(ss[1], 10, 64)
			if err != nil {
				return utils.WrapError(err.Error(), "invalid fields for cpuacct.stat content")
			}
			curCount += ticks
		}

		if i > 0 {
			// realtime usage := delta of cumulative ticks / ticks of one interval
			usedMilliCores := float64(curCount-count) / curTs.Sub(ts).Seconds() * jiffies * 1000.0
			log.Printf("cgroup %s cpu usage is %v (milli-cores)", cgroupPath, usedMilliCores)
		}

		ts = curTs
		count = curCount
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	return nil
}
