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

package optlyplugins

/// [Sohail]: Call it log file handler. And it should be generic. Called openFile/closeFile.
/// Use separate struct and put all those method in there. struct should have only three attributes
/// name, flag, filemode
import (
	"log"
	"os"
)

var logfile *os.File

// OpenLogFile opens the logfile
func OpenLogFile() {
	if logfile != nil {
		return
	}
	var err error
	// define as constant.
	logfile, err = os.OpenFile("/tmp/dat1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// CloseLogFile closes the logfile
func CloseLogFile() {
	if logfile == nil {
		return
	}
	if err := logfile.Close(); err != nil {
		log.Fatal(err)
	}
}

// WriteToLogFile writes provided string value to the logfile
func WriteToLogFile(str string) {
	if logfile == nil {
		return
	}
	if _, err := logfile.Write([]byte(str)); err != nil {
		log.Fatal(err)
	}
}
