// Copyright 2015 CoreStore Authors
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

package store_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store"
	"github.com/stretchr/testify/assert"
)

// These constants are here on purpose hard coded
func TestGetTable(t *testing.T) {

	tests := []struct {
		ti    csdb.Index
		isErr bool
	}{
		{ti: store.TableGroup, isErr: false},
		{ti: store.TableStore, isErr: false},
		{ti: store.TableWebsite, isErr: false},
		{ti: store.TableZMax, isErr: true},
	}

	for _, test := range tests {
		ts, err := store.GetTableStructure(test.ti)
		tn := store.GetTableName(test.ti)
		if test.isErr == false {
			assert.NoError(t, err)
			assert.NotNil(t, ts)
			assert.True(t, len(tn) > 1)
		} else {
			assert.Error(t, err)
			assert.Nil(t, ts)
			assert.Len(t, tn, 0)
		}
	}
}