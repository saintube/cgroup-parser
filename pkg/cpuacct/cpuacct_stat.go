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
	"log"
	"path/filepath"

	"github.com/saintube/cgroup-parser/pkg/utils"
)

const StatFileName = "cpuacct.stat"

func ErrInvalidStatArgs(argName string, errStr string) error {
	return utils.WrapError(errStr, "invalid argument "+argName+" for parsing cpuacct.stat")
}

func Stat(cgroupPath string, interval int, duration int) error {
	// workflow
	// 1. validate arguments
	// 2. get system Jifies
	// 3. loop reading cpuacct.stat and calculating the usage by intervals
	err := checkStatArgs(cgroupPath, interval, duration)
	if err != nil {
		return err
	}

	log.Println("start parse cpuacct.stat")

	return nil
}

func checkStatArgs(cgroupPath string, interval int, duration int) error {
	err := utils.FilePathExists(filepath.Join(cgroupPath, StatFileName))
	if err != nil {
		return ErrInvalidStatArgs("cgroup-path", err.Error())
	}
	if interval <= 0 {
		return ErrInvalidStatArgs("interval", "value should be larger than zero")
	}
	if duration <= 0 {
		return ErrInvalidStatArgs("interval", "value should be larger than zero")
	}
	return nil
}
