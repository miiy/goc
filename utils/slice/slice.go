package slice

import (
	"github.com/shopspring/decimal"
	"reflect"
)

// StructToSlice struct to slice
func StructToSlice(v interface{}) []interface{} {
	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}
	var s []interface{}
	for i := 0; i < to.NumField(); i++ {
		s = append(s, vo.Field(i).Interface())
	}
	return s
}

// StructToSliceByTagValues struct to slice by struct tag value
func StructToSliceByTagValues(f interface{}, tagName string, tagValues []string) []interface{} {
	to := reflect.TypeOf(f)
	vo := reflect.ValueOf(f)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}

	var s []interface{}
	for _, t := range tagValues {
		find := false
		for i := 0; i < to.NumField(); i++ {
			if t == to.Field(i).Tag.Get(tagName) {
				find = true
				s = append(s, vo.Field(i).Interface())
				break
			}
		}
		if !find {
			panic("not find tag " + t)
		}
	}

	return s
}

// SlicePluckInt pluck int column
func SlicePluckInt(v interface{}, fieldName string) []int {
	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}

	var r []int

	if to.Kind() == reflect.Slice {
		for i := 0; i < vo.Len(); i++ {
			ivo := vo.Index(i)
			if ivo.Kind() == reflect.Ptr {
				ivo = ivo.Elem()
			}
			r = append(r, int(ivo.FieldByName(fieldName).Int()))
			continue
		}
	}
	return r
}

// SlicePluckString pluck string column
func SlicePluckString(v interface{}, fieldName string) []string {
	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}

	var r []string

	if to.Kind() == reflect.Slice {
		for i := 0; i < vo.Len(); i++ {
			ivo := vo.Index(i)
			if ivo.Kind() == reflect.Ptr {
				ivo = ivo.Elem()
			}
			r = append(r, ivo.FieldByName(fieldName).String())
			continue
		}
	}
	return r
}

// SliceColumnSumDecimal column sum decimal
func SliceColumnSumDecimal(v interface{}, fieldName string) decimal.Decimal {
	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}

	var r decimal.Decimal

	if to.Kind() == reflect.Slice {
		for i := 0; i < vo.Len(); i++ {
			ivo := vo.Index(i)
			if ivo.Kind() == reflect.Ptr {
				ivo = ivo.Elem()
			}
			r = r.Add(ivo.FieldByName(fieldName).Interface().(decimal.Decimal))
		}
	}
	return r
}

// SliceGroupBy group by a slice column
// v slice
// fieldName is slice field name
// dist a map slice, such as map[string][]T
func SliceGroupBy(v interface{}, fieldName string, dist interface{}) {
	to := reflect.TypeOf(v)
	vo := reflect.ValueOf(v)
	if to.Kind() != reflect.Slice {
		panic("v must is a slice")
	}

	dto := reflect.TypeOf(dist)
	dvo := reflect.ValueOf(dist)
	if dto.Kind() != reflect.Map {
		panic("dist must is a map")
	}

	for i := 0; i < vo.Len(); i++ {
		iv := vo.Index(i)
		ivKey := iv.FieldByName(fieldName)
		dSlice := dvo.MapIndex(ivKey)
		if !dSlice.IsValid() {
			dSlice = reflect.MakeSlice(to, 0, 0)
		}
		dvo.SetMapIndex(ivKey, reflect.Append(dSlice, iv))
	}
}

// SliceColumnUniqueString string type column unique
func SliceColumnUniqueString(fieldName string, vs ...interface{}) []string {
	m := make(map[string]struct{})

	for _, v := range vs {
		to := reflect.TypeOf(v)
		vo := reflect.ValueOf(v)
		if to.Kind() == reflect.Ptr {
			to = to.Elem()
			vo = vo.Elem()
		}

		if to.Kind() == reflect.Slice {
			for i := 0; i < vo.Len(); i++ {
				ivo := vo.Index(i)
				if ivo.Kind() == reflect.Ptr {
					ivo = ivo.Elem()
				}
				ivof := ivo.FieldByName(fieldName)
				if ivof.Kind() != reflect.String {
					panic("field mus a string")
				}
				m[ivof.String()] = struct{}{}
			}
		}
	}

	var r []string
	for s := range m {
		r = append(r, s)
	}
	return r
}

// SliceColumnUniqueInt64 int64 type column unique
func SliceColumnUniqueInt64(fieldName string, vs ...interface{}) []int {
	m := make(map[int]struct{})

	for _, v := range vs {
		to := reflect.TypeOf(v)
		vo := reflect.ValueOf(v)
		if to.Kind() == reflect.Ptr {
			to = to.Elem()
			vo = vo.Elem()
		}

		if to.Kind() == reflect.Slice {
			for i := 0; i < vo.Len(); i++ {
				ivo := vo.Index(i)
				if ivo.Kind() == reflect.Ptr {
					ivo = ivo.Elem()
				}
				ivof := ivo.FieldByName(fieldName)
				if ivof.Kind() != reflect.Int64 {
					panic("field mus a int64")
				}
				m[int(ivof.Int())] = struct{}{}
			}
		}
	}

	var r []int
	for s := range m {
		r = append(r, s)
	}
	return r
}

// SliceElemPos slice elem pos
func SliceElemPos(needle interface{}, src interface{}) int {
	to := reflect.TypeOf(src)
	vo := reflect.ValueOf(src)
	if to.Kind() == reflect.Ptr {
		to = to.Elem()
		vo = vo.Elem()
	}
	if to.Kind() != reflect.Slice {
		panic("src must a slice")
	}
	for i := 0; i < vo.Len(); i++ {
		if reflect.ValueOf(needle).Interface() == vo.Index(i).Interface() {
			return i
		}
	}
	return 0
}
