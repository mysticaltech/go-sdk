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

// Package logging //
package logging

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var levelMap = map[LogLevel]zerolog.Level{
	LogLevelDebug:   zerolog.DebugLevel,
	LogLevelInfo:    zerolog.InfoLevel,
	LogLevelWarning: zerolog.WarnLevel,
	LogLevelError:   zerolog.ErrorLevel,
}

// FilteredLevelLogConsumer is an implementation of the OptimizelyLogConsumer that filters by log level
type FilteredLevelLogConsumer struct {
	logger *zerolog.Logger
}

// Log logs the message if it's log level is higher than or equal to the logger's set level
func (l *FilteredLevelLogConsumer) Log(level LogLevel, message string, fields map[string]interface{}) {
	l.logger.WithLevel(levelMap[level]).Fields(fields).Msg("[Optimizely] " + message)
}

// SetLogLevel changes the log level to the given level
func (l *FilteredLevelLogConsumer) SetLogLevel(level LogLevel) {
	childLogger := l.logger.Level(levelMap[level])
	l.logger = &childLogger
}

// NewFilteredLevelLogConsumer returns a new logger that logs to stdout
func NewFilteredLevelLogConsumer(level LogLevel, out io.Writer) *FilteredLevelLogConsumer {
	zerolog.SetGlobalLevel(levelMap[level])

	logger := log.Logger.Level(levelMap[level])
	return &FilteredLevelLogConsumer{
		logger: &logger,
	}
}
