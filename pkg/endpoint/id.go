//
// Copyright 2016 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package endpoint

import (
	"fmt"
	"strconv"
	"strings"
)

type PrefixType string

func (s PrefixType) String() string { return string(s) }

const (
	CiliumLocalIdPrefix  PrefixType = "cilium-local"
	CiliumGlobalIdPrefix            = "cilium-global"
	ContainerIdPrefix               = "container-id"
	DockerEndpointPrefix            = "docker-endpoint"
)

func NewCiliumID(id int64) string {
	return fmt.Sprintf("%s:%d", CiliumLocalIdPrefix, id)
}

func NewID(prefix PrefixType, id string) string {
	return string(prefix) + ":" + id
}

// Splits ID into prefix and id. No validation is performed on prefix
func SplitID(id string) (PrefixType, string) {
	if s := strings.Split(id, ":"); len(s) == 2 {
		return PrefixType(s[0]), s[1]
	} else {
		// default prefix
		return CiliumLocalIdPrefix, id
	}
}

// Parses id as cilium endpoint id and returns numeric portion
func ParseCiliumID(id string) (int64, error) {
	if prefix, id := SplitID(id); prefix != CiliumLocalIdPrefix {
		return 0, fmt.Errorf("not a cilium identifier")
	} else {
		if n, err := strconv.ParseInt(id, 0, 64); err != nil {
			return 0, fmt.Errorf("invalid numeric cilium id: %s", err)
		} else {
			return n, nil
		}
	}
}

// FIXME:
//  - Add docker ID and docker endpoint parsers

// Parses specified id and returns normalized id as string
func ParseID(id string) (PrefixType, string, error) {
	prefix, eid := SplitID(id)
	switch prefix {
	case CiliumLocalIdPrefix:
		if _, err := ParseCiliumID(id); err != nil {
			return "", "", err
		}
		return prefix, eid, nil
	case CiliumGlobalIdPrefix, ContainerIdPrefix, DockerEndpointPrefix:
		// FIXME: Validate IDs
		return prefix, eid, nil
	}

	return "", "", fmt.Errorf("unknown endpoint ID prefix \"%s\"", prefix)
}

// Parses specified id and returns normalized id as string
func ValidateID(id string) (PrefixType, string, error) {
	if prefix, _, err := ParseID(id); err != nil {
		return "", "", err
	} else {
		return prefix, id, nil
	}
}
