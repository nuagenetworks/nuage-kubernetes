package ovsdb

import (
	"fmt"
	"strings"
)

// NuageTableRow is an interface that all of the Nuage OVSDB table rows implement
type NuageTableRow interface {
	Equals(otherRow interface{}) bool
	CreateOVSDBRow(row map[string]interface{}) error
}

// UnMarshallOVSStringSet unmarshals a ovsdb column which is an array of strings
func UnMarshallOVSStringSet(data interface{}) ([]string, error) {
	var values []string
	var err error

	if datum, ok := data.(string); ok {
		values = append(values, datum)
	} else {
		var set []interface{}
		set, ok := data.([]interface{})
		if !ok {
			return values, fmt.Errorf("Invalid data")
		}

		if len(set) == 1 {
			values = append(values, (set[0]).(string))
		} else {
			var key string
			var ok bool
			if key, ok = set[0].(string); !ok {
				return nil, fmt.Errorf("Invalid type %+v", set)
			}

			if strings.Compare(key, "set") != 0 {
				return nil, fmt.Errorf("Invalid keyword %s", key)
			}

			valArr := (set[1]).([]interface{})
			for _, val := range valArr {
				values = append(values, val.(string))
			}
		}
	}

	return values, err
}
