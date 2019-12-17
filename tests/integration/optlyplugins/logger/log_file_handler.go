/****************************************************************************
 * Copyright 2019, Optimizely, Inc. and contributors                        *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

package logger

import (
	"log"
	"os"
	"sync"
)

// LogFileHandler manager class for log file handling
type LogFileHandler struct {
	Name      string
	Flag      int
	Mode      os.FileMode
	logfile   *os.File
	waitgroup sync.WaitGroup
}

// OpenFile opens the file with the provided name
func (l *LogFileHandler) OpenFile() {
	if l.logfile != nil {
		// Close already opened file
		l.CloseFile()
	}
	var err error
	l.waitgroup = sync.WaitGroup{}
	l.logfile, err = os.OpenFile(l.Name, l.Flag, l.Mode)
	if err != nil {
		log.Fatal(err)
	}
}

// CloseFile closes the file
func (l *LogFileHandler) CloseFile() {
	if l.logfile == nil {
		return
	}
	l.waitgroup.Wait()
	if err := l.logfile.Close(); err != nil {
		log.Fatal(err)
	}
	l.logfile = nil
}

// WriteToFile writes provided string value to the opened file
func (l *LogFileHandler) WriteToFile(str string) {
	if l.logfile == nil {
		return
	}
	l.waitgroup.Add(1)
	go func() {
		defer l.waitgroup.Done()
		if _, err := l.logfile.Write([]byte(str)); err != nil {
			log.Fatal(err)
		}
	}()
}
