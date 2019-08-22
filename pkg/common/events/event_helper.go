/*
Copyright 2019 Cloudera, Inc.  All rights reserved.

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

package events

import (
	"fmt"
)

func GetEventArgsAsStrings(result []string, generic []interface{}) error {
	if generic != nil && len(generic) > 0 {
		if result == nil || cap(result) != len(generic) {
			return fmt.Errorf("invalid length of arguments")
		}
		for idx, argument := range generic {
			argStr, ok := argument.(string)
			if ok {
				result[idx] = argStr
			}
		}
	}
	return nil
}