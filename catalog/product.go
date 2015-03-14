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

package catalog

import "github.com/corestoreio/csfw/eav"

type (
	ProductModel struct {
		// TBD
	}
)

var (
	_ eav.EntityTypeModeller                  = (*ProductModel)(nil)
	_ eav.EntityTypeTabler                    = (*ProductModel)(nil)
	_ eav.EntityTypeAttributeModeller         = (*ProductModel)(nil)
	_ eav.EntityTypeAdditionalAttributeTabler = (*ProductModel)(nil)
	_ eav.EntityAttributeCollectioner         = (*ProductModel)(nil)
)

func (c *ProductModel) TBD() {

}
func (c *ProductModel) AttributeCollection() {

}

func (c *ProductModel) TableNameBase() string {
	return GetTableName(TableProductEntity)
}

func (c *ProductModel) TableNameValue(i eav.ValueIndex) string {
	switch i {
	case eav.EntityTypeDateTime:
		return GetTableName(TableProductEntityDatetime)
	case eav.EntityTypeDecimal:
		return GetTableName(TableProductEntityDecimal)
	case eav.EntityTypeInt:
		return GetTableName(TableProductEntityInt)
	case eav.EntityTypeText:
		return GetTableName(TableProductEntityText)
	case eav.EntityTypeVarchar:
		return GetTableName(TableProductEntityVarchar)
	}
	return ""
}

func (c *ProductModel) TableNameAdditionalAttribute() string {
	return GetTableName(TableEavAttribute)
}

func Product() *ProductModel {
	return &ProductModel{}
}