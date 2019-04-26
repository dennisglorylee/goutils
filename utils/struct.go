package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
)

//dest 必须是指针
func CopyStruct(src, dest interface{}) error {
	ptrDestType := reflect.TypeOf(dest)
	if ptrDestType.Kind() != reflect.Ptr {
		return errors.New("dest type must be ptr")
	}

	destType := ptrDestType.Elem()

	srcValue := reflect.Indirect(reflect.ValueOf(src))
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	for i := 0; i < destType.NumField(); i++ {
		dstField := destType.Field(i)

		srcField := srcValue.FieldByName(dstField.Name)
		if srcField.IsValid() && srcField.Type() == dstField.Type {
			dstFieldValue := destValue.FieldByIndex(dstField.Index)
			if dstFieldValue.CanSet() {
				dstFieldValue.Set(srcField)
			}
		}
	}

	return nil
}

//dest 必须是数组的指针
func CopyStructs(src, dest interface{}, f func(interface{}, reflect.Type) interface{}) error {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Ptr && srcType.Kind() != reflect.Slice {
		return errors.New("src type must be slice or a slice address")
	}

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice {
		return errors.New("dest type must be a slice address")
	}

	ptrDestValue := reflect.ValueOf(dest)
	destValue := reflect.Indirect(ptrDestValue) //.Elem()
	destValue = destValue.Slice(0, destValue.Cap())

	srcValue := reflect.Indirect(reflect.ValueOf(src))

	destElemType := destType.Elem().Elem() //destType is slice address

	var isSliceElemPtr = false
	if destElemType.Kind() == reflect.Ptr {
		destElemType = destElemType.Elem()
		isSliceElemPtr = true
	}

	for i := 0; i < srcValue.Len(); i++ {
		srcElemValue := reflect.Indirect(srcValue.Index(i))

		destElemValuePtr := reflect.New(destElemType)
		destElemValue := reflect.Indirect(destElemValuePtr)

		for j := 0; j < destElemType.NumField(); j++ {
			dstElemTypeField := destElemType.Field(j)

			srcElemField := srcElemValue.FieldByName(dstElemTypeField.Name) //TODO cache
			if srcElemField.IsValid() {
				dstFieldValue := destElemValue.FieldByIndex(dstElemTypeField.Index)
				if dstFieldValue.CanSet() {
					if srcElemField.Type() == dstElemTypeField.Type {
						dstFieldValue.Set(srcElemField)
					} else {
						if f != nil {
							convSrcElemField := f(srcElemField.Interface(), dstElemTypeField.Type)
							if convSrcElemField != nil {
								dstFieldValue.Set(reflect.ValueOf(convSrcElemField))
							}
						}
					}
				}
			}
		}

		if isSliceElemPtr {
			destValue = reflect.Append(destValue, destElemValuePtr)
		} else {
			destValue = reflect.Append(destValue, destElemValue)
		}
	}

	ptrDestValue.Elem().Set(destValue.Slice(0, destValue.Len()))
	return nil
}

func AutoGenTags(vo interface{}, tagDefs map[string]TAG_STYLE) string {
	voType := reflect.TypeOf(vo)
	if voType.Kind() == reflect.Ptr {
		voType = voType.Elem()
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("type %s struct{\n", voType.Name()))

	var sortedTagName []string
	for k := range tagDefs {
		sortedTagName = append(sortedTagName, k)
	}
	sort.Strings(sortedTagName)

	for i := 0; i < voType.NumField(); i++ {
		dstField := voType.Field(i)

		var tagBuffer bytes.Buffer
		for _, tagName := range sortedTagName {
			tagStr := ConvertFieldStyle(dstField.Name, tagDefs[tagName])
			if tagStr != "" {
				tagBuffer.WriteString(fmt.Sprintf(`%s:"%s" `, tagName, tagStr))
			}
		}

		buffer.WriteString(fmt.Sprintf("%s %s `%s`\n", dstField.Name,
			dstField.Type.String(), tagBuffer.String())) //pkgStr+dstField.Type.Name()
	}

	buffer.WriteString("}")

	return buffer.String()
}
