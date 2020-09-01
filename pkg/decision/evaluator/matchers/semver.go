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

// Package matchers //
package matchers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/optimizely/go-sdk/pkg/logging"

	"github.com/optimizely/go-sdk/pkg/decision/reasons"
	"github.com/optimizely/go-sdk/pkg/entities"

	"github.com/pkg/errors"
)

const (
	buildSeperator      = "+"
	preReleaseSeperator = "-"
	whiteSpace          = " "
)

var digitCheck = regexp.MustCompile(`^[0-9]+$`)

// SemanticVersion defines the class
type SemanticVersion struct {
	Condition string // condition is always a string here
}

func (sv SemanticVersion) compareVersion(attribute string) (int, error) {

	targetedVersionParts, err := sv.splitSemanticVersion(sv.Condition)
	if err != nil {
		return 0, err
	}
	versionParts, e := sv.splitSemanticVersion(attribute)
	if e != nil {
		return 0, e
	}

	// Up to the precision of targetedVersion, expect version to match exactly.
	for idx := range targetedVersionParts {

		switch {
		case len(versionParts) <= idx:
			if sv.isPreReleaseOrBuild(attribute) == 1 {
				return 1, nil
			}
			return -1, nil
		case !sv.isNumber(versionParts[idx]):
			// Compare strings
			if versionParts[idx] < targetedVersionParts[idx] {
				return -1, nil
			} else if versionParts[idx] > targetedVersionParts[idx] {
				return 1, nil
			}
		case sv.isNumber(targetedVersionParts[idx]): // both targetedVersionParts and versionParts are digits
			if sv.toInt(versionParts[idx]) < sv.toInt(targetedVersionParts[idx]) {
				return -1, nil
			} else if sv.toInt(versionParts[idx]) > sv.toInt(targetedVersionParts[idx]) {
				return 1, nil
			}
		default:
			return -1, nil
		}
	}

	if sv.isPreReleaseOrBuild(attribute) == 1 && sv.isPreReleaseOrBuild(sv.Condition) != 1 {
		return -1, nil
	}

	return 0, nil
}

func (sv SemanticVersion) splitSemanticVersion(targetedVersion string) ([]string, error) {

	if sv.hasWhiteSpace(targetedVersion) {
		return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
	}

	targetPrefix := targetedVersion
	var targetSuffix string

	if sv.isPreReleaseOrBuild(targetedVersion) != 0 {
		// More than one occurrence of build separator not allowed
		if strings.Count(targetedVersion, buildSeperator) > 1 {
			return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
		}

		sep := buildSeperator
		if sv.isPreReleaseOrBuild(targetedVersion) == 1 {
			sep = preReleaseSeperator
		}
		// this is going to slit with the first occurrence.
		index := strings.Index(targetedVersion, sep)
		targetPrefix = targetedVersion[:index]
		targetSuffix = targetedVersion[index+1:]

		// both prefix and suffix should be present in case a separator is present
		if targetPrefix == "" || targetSuffix == "" {
			return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
		}
		// dont compare build meta-data
		if tSuffix, err := sv.getValidTargetSuffix(sep, targetSuffix); err == nil {
			targetSuffix = tSuffix
		} else {
			return []string{}, err
		}
	}

	// prerelease must only comprise only ASCII alphanumerics and hyphens
	if !sv.isValidPreReleaseOrBuild(targetSuffix) {
		return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
	}

	// Expect a version string of the form x.y.z
	// split all dots with SplitAfter
	targetedVersionParts := strings.Split(targetPrefix, ".")

	targetVersionPartsCount := len(targetedVersionParts)
	if targetVersionPartsCount > 3 || targetVersionPartsCount == 0 {
		return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
	}

	for i := 0; i < len(targetedVersionParts); i++ {
		if !sv.isNumber(targetedVersionParts[i]) {
			return []string{}, errors.New(string(reasons.AttributeFormatInvalid))
		}
	}

	if targetSuffix != "" {
		targetedVersionParts = append(targetedVersionParts, targetSuffix)
	}
	return targetedVersionParts, nil
}

func (sv SemanticVersion) getValidTargetSuffix(separator, str string) (string, error) {
	if separator == buildSeperator {
		// build must only comprise only ASCII alphanumerics and hyphens
		if !sv.isValidPreReleaseOrBuild(str) {
			return "", errors.New(string(reasons.AttributeFormatInvalid))
		}
		return "", nil
	} else if index := strings.Index(str, buildSeperator); index != -1 {
		metaData := str[index:]
		// build must only comprise only ASCII alphanumerics and hyphens
		if !sv.isValidPreReleaseOrBuild(metaData[1:]) {
			return "", errors.New(string(reasons.AttributeFormatInvalid))
		}
		return strings.ReplaceAll(str, metaData, ""), nil
	}
	return str, nil
}

func (sv SemanticVersion) isNumber(str string) bool {
	return digitCheck.MatchString(str)
}

func (sv SemanticVersion) toInt(str string) int {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return i
}

// Returns true if string only comprises only ASCII alphanumerics and hyphens
func (sv SemanticVersion) isValidPreReleaseOrBuild(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9-]*$`).MatchString(str)
}

// Returns 1 if prerelease, -1 if build, 0 if neither
func (sv SemanticVersion) isPreReleaseOrBuild(str string) int {
	if strings.Contains(str, preReleaseSeperator) || strings.Contains(str, buildSeperator) {
		indexBuildSeparator := strings.Index(str, buildSeperator)
		indexPreReleaseSeparator := strings.Index(str, preReleaseSeperator)
		switch {
		case indexBuildSeparator == -1:
			return 1
		case indexPreReleaseSeparator == -1:
			return -1
		default:
			if indexBuildSeparator > indexPreReleaseSeparator {
				return 1
			}
			return -1
		}
	}
	return 0
}

func (sv SemanticVersion) hasWhiteSpace(str string) bool {
	return str == "" || strings.Contains(str, whiteSpace)
}

// SemverEvaluator is a help function to wrap a common evaluation code
func SemverEvaluator(cond entities.Condition, user entities.UserContext) (int, error) {

	if stringValue, ok := cond.Value.(string); ok {
		attributeValue, err := user.GetStringAttribute(cond.Name)
		if err != nil {
			return 0, err
		}
		semVer := SemanticVersion{stringValue}
		comparison, e := semVer.compareVersion(attributeValue)
		if e != nil {
			return 0, e
		}
		return comparison, nil
	}
	return 0, fmt.Errorf("audience condition %s evaluated to NULL because the condition value type is not supported", cond.Name)
}

// SemverEqMatcher returns true if the user's semver attribute is equal to the semver condition value
func SemverEqMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	comparison, err := SemverEvaluator(condition, user)
	if err != nil {
		return false, err
	}
	return comparison == 0, nil
}

// SemverGeMatcher returns true if the user's semver attribute is greater or equal to the semver condition value
func SemverGeMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	comparison, err := SemverEvaluator(condition, user)
	if err != nil {
		return false, err
	}
	return comparison >= 0, nil
}

// SemverGtMatcher returns true if the user's semver attribute is greater than the semver condition value
func SemverGtMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	comparison, err := SemverEvaluator(condition, user)
	if err != nil {
		return false, err
	}
	return comparison > 0, nil
}

// SemverLeMatcher returns true if the user's semver attribute is less than or equal to the semver condition value
func SemverLeMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	comparison, err := SemverEvaluator(condition, user)
	if err != nil {
		return false, err
	}
	return comparison <= 0, nil
}

// SemverLtMatcher returns true if the user's semver attribute is less than the semver condition value
func SemverLtMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	comparison, err := SemverEvaluator(condition, user)
	if err != nil {
		return false, err
	}
	return comparison < 0, nil
}
