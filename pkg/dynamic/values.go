package dynamic

import (
	"github.com/vektah/gqlgen/neelance/common"
	"time"
	"github.com/vektah/gqlgen/neelance/schema"
	"github.com/vektah/gqlgen/graphql"
)

type Value interface {
	common.Type
	Type() common.Type
	Marshaller() graphql.Marshaler
	//TODO:
	//Validate() error
}

// enforce interface
var (
	_ Value = &Object{}
	_ Value = &Array{}
	_ Value = &String{}
	_ Value = &Enum{}
	_ Value = &Null{}
	_ Value = &Float{}
	_ Value = &Int{}
	_ Value = &Time{}
)

type Object struct {
	*schema.Object
	Data *OrderedMap
}

type Array struct {
	*common.List
	Data []Value
}

type Int struct {
	*schema.Scalar
	Data int
}

type String struct {
	*schema.Scalar
	Data string
}

type Enum struct {
	*schema.Enum
	Data string
}

type Float struct {
	*schema.Scalar
	Data float64
}

type Bool struct {
	*schema.Scalar
	Data bool
}

type Time struct {
	*schema.Scalar
	Data time.Time
}

type Null struct{}

func (t *Null) Kind() string   { return "NULL" }
func (t *Null) String() string { return "null" }


func (t *Object) Type() common.Type {
	return t.Object
}
func (t *Array) Type() common.Type {
	return t.List
}
func (t *Int) Type() common.Type {
	return t.Scalar
}
func (t *Float) Type() common.Type {
	return t.Scalar
}
func (t *String) Type() common.Type {
	return t.Scalar
}
func (t *Enum) Type() common.Type {
	return t.Enum
}
func (t *Bool) Type() common.Type {
	return t.Scalar
}
func (t *Time) Type() common.Type {
	return t.Scalar
}
func (t *Null) Type() common.Type {
	return nil
}

func (t *Object) Marshaller() graphql.Marshaler {
	items := t.Data.Items()
	fieldMap := graphql.NewOrderedMap(len(items))
	for i, item := range items {
		fieldMap.Keys[i] = item.Key
		fieldMap.Values[i] = item.Value.Marshaller()
	}
	return fieldMap
}
func (t *Array) Marshaller() graphql.Marshaler {
	var array graphql.Array
	for _, val := range t.Data {
		array = append(array, val.Marshaller())
	}
	return array
}
func (t *Int) Marshaller() graphql.Marshaler {
	return graphql.MarshalInt(t.Data)
}
func (t *Float) Marshaller() graphql.Marshaler {
	return graphql.MarshalFloat(t.Data)
}
func (t *String) Marshaller() graphql.Marshaler {
	return graphql.MarshalString(t.Data)
}
func (t *Enum) Marshaller() graphql.Marshaler {
	return graphql.MarshalString(t.Data)
}
func (t *Bool) Marshaller() graphql.Marshaler {
	return graphql.MarshalBoolean(t.Data)
}
func (t *Time) Marshaller() graphql.Marshaler {
	return graphql.MarshalTime(t.Data)
}
func (t *Null) Marshaller() graphql.Marshaler {
	return graphql.Null
}

// preserving order matters
type OrderedMap struct {
	Keys   []string
	Values []Value
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{}
}

func (m *OrderedMap) Set(key string, value Value) {
	for i, k := range m.Keys {
		if key == k {
			m.Values[i] = value
			return
		}
	}
	m.Keys = append(m.Keys, key)
	m.Values = append(m.Values, value)
}

func (m *OrderedMap) Get(key string) Value {
	for i, k := range m.Keys {
		if key == k {
			return m.Values[i]
		}
	}
	return nil
}

func (m *OrderedMap) Delete(key string) {
	for i, k := range m.Keys {
		if key == k {
			m.Keys = append(m.Keys[:i], m.Keys[i+1:]...)
			m.Values = append(m.Values[:i], m.Values[i+1:]...)
			return
		}
	}
}

func (m *OrderedMap) Items() []struct {
	Key   string
	Value Value
} {
	var items []struct {
		Key   string;
		Value Value
	}
	for i, k := range m.Keys {
		items = append(items, struct {
			Key   string;
			Value Value
		}{
			Key:   k,
			Value: m.Values[i],
		})
	}
	return items
}
