// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dosa

import (
	"testing"

	"time"

	"io"

	"github.com/stretchr/testify/assert"
)

type SinglePrimaryKeyNoParen struct {
	Entity     `dosa:"primaryKey=PrimaryKey"`
	PrimaryKey int64
	Data       string
}

// happy path: A single primaryKey becomes the partition key
func TestSinglePrimaryKeyNoParen(t *testing.T) {
	dosaTable, err := TableFromInstance(&SinglePrimaryKeyNoParen{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"primarykey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, 0, len(dosaTable.Key.ClusteringKeys))
}

func TestNilPointer(t *testing.T) {
	dosaTable, err := TableFromInstance((*SinglePrimaryKeyNoParen)(nil))
	assert.Nil(t, err)
	assert.Equal(t, []string{"primarykey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, 0, len(dosaTable.Key.ClusteringKeys))
}

type SinglePrimaryKey struct {
	Entity     `dosa:"primaryKey=(PrimaryKey)"`
	PrimaryKey int64
	Data       string
}

// happy path: A single primaryKey becomes the partition key
func TestSinglePrimaryKey(t *testing.T) {
	dosaTable, err := TableFromInstance(&SinglePrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"primarykey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, 0, len(dosaTable.Key.ClusteringKeys))
}

func BenchmarkSingleKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TableFromInstance(&SinglePrimaryKey{})
	}
}

type SinglePartitionKey struct {
	Entity     `dosa:"primaryKey=PrimaryKey"`
	PrimaryKey int64
	data       string
}

func TestSinglePartitionKey(t *testing.T) {
	dosaTable, err := TableFromInstance(&SinglePartitionKey{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"primarykey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, 0, len(dosaTable.Key.ClusteringKeys))
}

// unhappy path: this struct doesn't have anything specified for pk
type NoPrimaryKey struct {
	Entity     `dosa:"primaryKey="`
	PrimaryKey int64
	data       string
}

// unhappy path: If there is no field marked with a primary nor partition key, throw an error
func TestNoPrimaryKey(t *testing.T) {
	dosaTable, err := TableFromInstance(&NoPrimaryKey{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
}

// unhappy path: this struct has an empty primary key
type EmptyPrimaryKey struct {
	Entity     `dosa:"primaryKey=()"`
	PrimaryKey int64
	data       string
}

// unhappy path: If there is no field marked with a primary nor partition key, throw an error
func TestEmptyPrimaryKey(t *testing.T) {
	dosaTable, err := TableFromInstance(&EmptyPrimaryKey{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
}

type PrimaryKeyWithSecondaryRange struct {
	Entity         `dosa:"primaryKey=(PartKey,PrimaryKey  )"`
	PartKey        int64
	PrimaryKey     int64
	data, moredata string
}

func TestPrimaryKeyWithSecondaryRange(t *testing.T) {
	dosaTable, err := TableFromInstance(&PrimaryKeyWithSecondaryRange{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"partkey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, []*ClusteringKey{{"primarykey", false}}, dosaTable.Key.ClusteringKeys)
}

type PrimaryKeyWithDescendingRange struct {
	Entity     `dosa:"primaryKey=(PartKey,PrimaryKey desc )"`
	PartKey    int64
	PrimaryKey int64
	data       string
}

func TestPrimaryKeyWithDescendingRange(t *testing.T) {
	dosaTable, err := TableFromInstance(&PrimaryKeyWithDescendingRange{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"partkey"}, dosaTable.Key.PartitionKeys)
	assert.Equal(t, []*ClusteringKey{{"primarykey", true}}, dosaTable.Key.ClusteringKeys)
}

type MultiComponentPrimaryKey struct {
	Entity         `dosa:"primaryKey=((PartKey, AnotherPartKey))"`
	PartKey        int64
	AnotherPartKey int64
	data           string
}

func TestMultiComponentPrimaryKey(t *testing.T) {
	dosaTable, err := TableFromInstance(&MultiComponentPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, []string{"partkey", "anotherpartkey"}, dosaTable.Key.PartitionKeys)
	assert.Nil(t, dosaTable.Key.ClusteringKeys)
}

type InvalidDosaAttribute struct {
	Entity `dosa:"oopsie, primaryKey=Oops"`
	Oops   int64
}

func TestInvalidDosaAttribute(t *testing.T) {
	dosaTable, err := TableFromInstance(&InvalidDosaAttribute{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "oopsie")
}

func TestStringify(t *testing.T) {
	dosaTable, _ := TableFromInstance(&SinglePrimaryKey{})
	assert.Contains(t, dosaTable.String(), dosaTable.Name)
	assert.Contains(t, dosaTable.String(), "primarykey")
}

type MissingCloseParen struct {
	Entity            `dosa:"primaryKey=(MissingCloseParen"`
	MissingCloseParen int64
}

func TestMissingCloseParen(t *testing.T) {
	dosaTable, err := TableFromInstance(&MissingCloseParen{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "MissingCloseParen")
}

type MissingAnnotation struct {
	Entity
	Oops int64
}

func TestMissingAnnotation(t *testing.T) {
	dosaTable, err := TableFromInstance(&MissingAnnotation{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "struct")
	assert.Contains(t, err.Error(), "tag")
}

type AllTypes struct {
	Entity     `dosa:"primaryKey=BoolType"`
	BoolType   bool
	Int32Type  int32
	Int64Type  int64
	DoubleType float64
	StringType string
	BlobType   []byte
	TimeType   time.Time
	UUIDType   UUID
}

func TestAllTypes(t *testing.T) {
	dosaTable, err := TableFromInstance(&AllTypes{})
	assert.NotNil(t, dosaTable)
	assert.Nil(t, err)
	cds := dosaTable.Columns
	assert.Len(t, cds, 8)
	for _, cd := range cds {
		name, err := NormalizeName(cd.Name)
		assert.NoError(t, err)
		switch name {
		case "booltype":
			assert.Equal(t, Bool, cd.Type)
		case "int32type":
			assert.Equal(t, Int32, cd.Type)
		case "int64type":
			assert.Equal(t, Int64, cd.Type)
		case "doubletype":
			assert.Equal(t, Double, cd.Type)
		case "stringtype":
			assert.Equal(t, String, cd.Type)
		case "blobtype":
			assert.Equal(t, Blob, cd.Type)
		case "timetype":
			assert.Equal(t, Timestamp, cd.Type)
		case "uuidtype":
			assert.Equal(t, TUUID, cd.Type)
		default:
			assert.Fail(t, "unexpected column name", name)
		}
	}
}

type UnsupportedType struct {
	Entity    `dosa:"primaryKey=BoolType"`
	BoolType  bool
	UnsupType float32
}

func TestUnsupportedType(t *testing.T) {
	dosaTable, err := TableFromInstance(&UnsupportedType{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "float32")
	assert.Contains(t, err.Error(), "UnsupType")
}

type KeyFieldNameTypo struct {
	Entity   `dosa:"primaryKey=BoolHype"`
	BoolType bool
}

func TestKeyFieldNameTypo(t *testing.T) {
	dosaTable, err := TableFromInstance(&KeyFieldNameTypo{})
	assert.Nil(t, dosaTable)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BoolHype")
}

type UnexportedFieldType struct {
	Entity          `dosa:"primaryKey=BoolType"`
	BoolType        bool
	unexported      string
	unsupportedType io.Reader
}

func TestIgnoreUnexportedFields(t *testing.T) {
	table, err := TableFromInstance(&UnexportedFieldType{})
	assert.NoError(t, err)
	assert.Len(t, table.ColToField, 1)
	assert.NotContains(t, table.ColToField, "unexported")
	assert.NotContains(t, table.ColToField, "unsupportedType")
	assert.Len(t, table.Columns, 1)
}

type IgnoreTagType struct {
	Entity          `dosa:"primaryKey=BoolType"`
	BoolType        bool
	Exported        string    `dosa:"-"`
	UnsupportedType io.Reader `dosa:"-"`
}

func TestIgnoreTag(t *testing.T) {
	table, err := TableFromInstance(&IgnoreTagType{})
	assert.NoError(t, err)
	assert.Len(t, table.ColToField, 1)
	// TODO: should be checking normalized names once normalization is done
	assert.NotContains(t, table.ColToField, "Exported")
	assert.NotContains(t, table.ColToField, "UnsupportedType")
	assert.Len(t, table.Columns, 1)
}

func TestExtraStuffInClusteringKeyDecl(t *testing.T) {
	type BadClusteringKeyDefinition struct {
		Entity     `dosa:"primaryKey=(BoolType,StringType asc asc)"`
		BoolType   bool
		StringType string
	}

	table, err := TableFromInstance(&BadClusteringKeyDefinition{})
	assert.Nil(t, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lustering")
}

func TestDuplicateKey(t *testing.T) {
	type DuplicateKeyDefinition struct {
		Entity     `dosa:"primaryKey=(BoolType,BoolType)"`
		BoolType   bool
		StringType string
	}

	table, err := TableFromInstance(&DuplicateKeyDefinition{})
	assert.Nil(t, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uplicate")
}

func TestInvalidStructName(t *testing.T) {
	type ăBădNăme struct {
		Entity   `dosa:"primaryKey=BoolType"`
		BoolType bool
	}
	table, err := TableFromInstance(&ăBădNăme{})
	assert.Nil(t, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ormalize")
}
func TestInvalidFieldInTag(t *testing.T) {
	type HasInvalidCharInTag struct {
		Entity   `dosa:"primaryKey=(ABădNăme)"`
		ABădNăme bool
	}
	table, err := TableFromInstance(&HasInvalidCharInTag{})
	assert.Nil(t, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

/*
 These tests do not currently pass, but I think they should
*/
func TestRenameToInvalidName(t *testing.T) {
	type InvalidRename struct {
		Entity    `dosa:"primaryKey=BoolType"`
		BoolType  bool
		AGoodName string `dosa:"name=ABădNăme"`
	}
	table, err := TableFromInstance(&InvalidRename{})
	assert.Nil(t, table)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid name tag: name=ABădNăme")
}

type BadColNameButRenamed struct {
	Entity   `dosa:"primaryKey=(ABădNăme)"`
	ABădNăme bool `dosa:"name=goodname"`
}

func TestRenameColumnToValidName(t *testing.T) {
	table, err := TableFromInstance(&BadColNameButRenamed{})
	assert.NoError(t, err)
	assert.Equal(t, []string{"goodname"}, table.Key.PartitionKeys)
}

type StructWithUnannotatedEntity struct {
	Entity `notdosa:"covers a rare test case"`
}

func TestRenameStructToValidName(t *testing.T) {
	type ABădNăme struct {
		Entity `dosa:"name=agoodname,primaryKey=Dummy"`
		Dummy  bool
	}
	table, err := TableFromInstance(&ABădNăme{})
	assert.NotNil(t, table)
	assert.NoError(t, err)
}
