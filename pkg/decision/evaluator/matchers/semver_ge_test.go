/****************************************************************************
 * Copyright 2020, Optimizely, Inc. and contributors                        *
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

package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/optimizely/go-sdk/pkg/entities"
)

func TestSemverGeMatcher(t *testing.T) {
	matcher := SemVerGeMatcher{
		Condition: entities.Condition{
			Match: "semver_ge",
			Value: "2.0",
			Name:  "version",
		},
	}

	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"version": "2.0.0",
		},
	}
	result, err := matcher.Match(user)
	assert.NoError(t, err)
	assert.True(t, result)

	user = entities.UserContext{
		Attributes: map[string]interface{}{
			"version": "2.9",
		},
	}

	result, err = matcher.Match(user)
	assert.NoError(t, err)
	assert.True(t, result)

	user = entities.UserContext{
		Attributes: map[string]interface{}{
			"version": "1.9",
		},
	}

	result, err = matcher.Match(user)
	assert.NoError(t, err)
	assert.False(t, result)

	// Test attribute not found
	user = entities.UserContext{
		Attributes: map[string]interface{}{
			"version1": "2.0",
		},
	}

	_, err = matcher.Match(user)
	assert.Error(t, err)
}