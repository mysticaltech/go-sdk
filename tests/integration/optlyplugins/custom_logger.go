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

import (
	"fmt"

	"github.com/optimizely/go-sdk/pkg/logging"
)

// CustomLogger custom implementation of OptimizelyLogConsumer interface
type CustomLogger struct {
	ScenarioID string
	level      logging.LogLevel
}

// Log stores the message in the logfile
func (l CustomLogger) Log(level logging.LogLevel, message string, fields map[string]interface{}) {
	if l.level <= level {
		// prepends the name and log level to the message
		//TODO: what is fields[name]
		message = fmt.Sprintf("[OPTIMIZELY][%s][Request: %s][%s] %s", level, l.ScenarioID, fields["name"], message)

		//TODO: Tis can be expensive, for now it's fine but we should call it async.
		WriteToLogFile(message)
	}
}

// SetLogLevel changes the log level to the given level
func (l CustomLogger) SetLogLevel(level logging.LogLevel) {
	l.level = level
}
