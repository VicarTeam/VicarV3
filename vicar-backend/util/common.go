package util

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/google/uuid"
)

func Ptr[T any](v T) *T {
	return &v
}

func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func Unptr[T any](val *T, fallback T) T {
	if val == nil {
		return fallback
	}

	return *val
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

func EqSlice[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, item := range a {
		if item != b[i] {
			return false
		}
	}

	return true
}

func ParseEmptyableUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}

	return Ptr(uuid.MustParse(s))
}

func UUIDToEmptyableString(u *uuid.UUID) string {
	if u == nil {
		return ""
	}

	return u.String()
}

func EmptyStringToPtr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func ContainsPredicate[T any](slice []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(slice, predicate)
}

func MapArray[T any, K any](slice []T, mapper func(T) K) []K {
	result := make([]K, len(slice))

	for i, item := range slice {
		result[i] = mapper(item)
	}

	return result
}

func CountArray[T any](slice []T, predicate func(T) bool) int {
	count := 0

	for _, item := range slice {
		if predicate(item) {
			count++
		}
	}

	return count
}

func AnyArray[T any](slice []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(slice, predicate)
}

func FilterArray[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

func FindArray[T any](slice []T, predicate func(T) bool) *T {
	for _, item := range slice {
		if predicate(item) {
			return &item
		}
	}

	return nil
}

func GetMapElementAt[TKey comparable, TValue any](m map[TKey]TValue, index int) (*TKey, *TValue) {
	i := 0
	for k, v := range m {
		if i == index {
			return &k, &v
		}

		i++
	}

	return nil, nil
}

func AnyToString(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", v)
}

func SliceToMap[T any, K comparable, V any](slice []T, keyMapper func(T) K, valueMapper func(T) V) map[K]V {
	result := map[K]V{}

	for _, item := range slice {
		result[keyMapper(item)] = valueMapper(item)
	}

	return result
}

func MapToSlice[T any, K comparable, V any](m map[K]V, mapper func(K, V) T) []T {
	result := make([]T, 0, len(m))

	for k, v := range m {
		result = append(result, mapper(k, v))
	}

	return result
}
