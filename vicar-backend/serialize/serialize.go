package serialize

import "github.com/gofiber/fiber/v2"

type Serializer[T any] interface {
	Serialize(obj T, args ...any) any
}

type SerializerWithFlags[T any] interface {
	Serializer[T]
	ParseFlags(c *fiber.Ctx)
}

func Do[T any](serializer Serializer[T], obj T, args ...any) any {
	return serializer.Serialize(obj, args...)
}

func DoNilable[T any](serializer Serializer[T], obj *T, args ...any) any {
	if obj == nil {
		return nil
	}

	return serializer.Serialize(*obj, args...)
}

func DoWithFlags[T any](serializer SerializerWithFlags[T], c *fiber.Ctx, obj T, args ...any) any {
	serializer.ParseFlags(c)
	return serializer.Serialize(obj, args...)
}

func DoArray[T any](serializer Serializer[T], objs []T, args ...any) []any {
	result := make([]any, len(objs))
	for i, obj := range objs {
		result[i] = Do[T](serializer, obj, args...)
	}
	return result
}

func DoFilteredArray[T any](serializer Serializer[T], objs []T, predicate func(T) bool, args ...any) []any {
	result := make([]any, 0, len(objs))
	for _, obj := range objs {
		if predicate(obj) {
			result = append(result, Do[T](serializer, obj, args...))
		}
	}
	return result
}

func DoPtrArray[T any](serializer Serializer[T], objs []*T, args ...any) []any {
	result := make([]any, 0, len(objs))
	for _, obj := range objs {
		if obj != nil {
			result = append(result, Do[T](serializer, *obj, args...))
		}
	}
	return result
}

func DoFilteredPtrArray[T any](serializer Serializer[T], objs []*T, predicate func(*T) bool, args ...any) []any {
	result := make([]any, 0, len(objs))
	for _, obj := range objs {
		if obj != nil {
			if predicate(obj) {
				result = append(result, Do[T](serializer, *obj, args...))
			}
		}
	}
	return result
}

func DoArrayWithFlags[T any](serializer SerializerWithFlags[T], c *fiber.Ctx, objs []T, args ...any) []any {
	result := make([]any, len(objs))
	serializer.ParseFlags(c)
	for i, obj := range objs {
		result[i] = Do[T](serializer, obj, args...)
	}

	return result
}

func DoVarargs[T any](serializer Serializer[T], objs ...T) []any {
	return DoArray[T](serializer, objs)
}

func DoVarargsWithFlags[T any](serializerWithFlags SerializerWithFlags[T], c *fiber.Ctx, objs ...T) []any {
	return DoArrayWithFlags[T](serializerWithFlags, c, objs)
}

func JSON[T any](c *fiber.Ctx, serializer Serializer[T], obj T, args ...any) error {
	return c.JSON(Do[T](serializer, obj, args...))
}

func JSONNilable[T any](c *fiber.Ctx, serializer Serializer[T], obj *T, args ...any) error {
	return c.JSON(DoNilable[T](serializer, obj, args...))
}

func JSONWithFlags[T any](c *fiber.Ctx, serializer SerializerWithFlags[T], obj T, args ...any) error {
	return c.JSON(DoWithFlags[T](serializer, c, obj, args...))
}

func JSONArray[T any](c *fiber.Ctx, serializer Serializer[T], objs []T, args ...any) error {
	return c.JSON(DoArray[T](serializer, objs, args...))
}

func JSONPtrArray[T any](c *fiber.Ctx, serializer Serializer[T], objs []*T, args ...any) error {
	return c.JSON(DoPtrArray[T](serializer, objs, args...))
}

func JSONArrayWithFlags[T any](c *fiber.Ctx, serializer SerializerWithFlags[T], objs []T, args ...any) error {
	return c.JSON(DoArrayWithFlags[T](serializer, c, objs, args...))
}

func JSONVarargs[T any](c *fiber.Ctx, serializer Serializer[T], objs ...T) error {
	return c.JSON(DoVarargs[T](serializer, objs...))
}

func JSONVarargsWithFlags[T any](c *fiber.Ctx, serializer SerializerWithFlags[T], objs ...T) error {
	return c.JSON(DoVarargsWithFlags(serializer, c, objs...))
}

type Context map[string]any

var NoCtx = Context(map[string]any{})

func Ctx(values map[string]any) Context {
	return Context(values)
}

func (c Context) Set(key string, value any) Context {
	c[key] = value
	return c
}

func (c Context) Get(key string) any {
	return c[key]
}

func (c Context) Has(key string) bool {
	_, ok := c[key]
	return ok
}
