# PHP Serialize for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/stlong5/phpserialize.svg)](https://pkg.go.dev/github.com/stlong5/phpserialize)

A robust Go package that provides PHP-compatible serialize (`Marshal`) and unserialize (`Unmarshal`) functions. This
library is designed for reliable data exchange between Go and PHP applications, particularly when dealing with
persistent data formats like session data or database-stored serialized values.

## Features✨

✅ **Full PHP Compatibility**: Supports PHP 4, 5, 7, and 8 serialization formats.  
✅ **All PHP Types**: Handles null, bool, int, float, string, array, and object.     
✅ **Security Options**: Includes nesting depth limits and allowed class filtering for robust handling of untrusted
data.    
✅ **Unified API**: A single Option interface is used to customize both Marshal and Unmarshal behavior.  
✅ **Special Type Handling**: Correctly manages NaN, Inf, -Inf, and PHP's private/protected property naming
conventions.  
✅ **Binary Safe**: Properly handles strings containing null bytes and non-UTF-8 data.   
✅ Zero Dependencies and Go 1.24+ support.

## Installation

```shell
go get github.com/stlong5/phpserialize
```

## Quick Start🚀

This example shows how to serialize a Go map and then unserialize it back.

```go
package main

import (
	"fmt"
	"log"
	"github.com/stlong5/phpserialize"
)

func main() {
	// Go data structure
	data := map[string]interface{}{
		"name":  "John Doe",
		"age":   30,
		"roles": []interface{}{"admin", "user"},
	}

	// 1. Marshal (Serialize) Go to PHP format
	serialized, err := phpserialize.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Serialized:", serialized)
	// Output: a:3:{s:4:"name";s:8:"John Doe";s:3:"age";i:30;s:5:"roles";a:2:{i:0;s:5:"admin";i:1;s:4:"user";}}

	// 2. Unmarshal (Unserialize) PHP format back to Go
	unserialized, err := phpserialize.Unmarshal(serialized)
	if err != nil {
		log.Fatal(err)
	}

	// The result is a map[string]interface{}
	fmt.Printf("Unserialized: %+v\n", unserialized)
	// Output: Unserialized: map[age:30 name:John Doe roles:[admin user]]
}
```

### PHP Objects

PHP objects are represented in Go using the `phpserialize.PHPObject` struct.

```go
// PHPObject definition
type PHPObject struct {
ClassName  string
Properties map[string]interface{}
}

// Serialization
phpObject := phpserialize.PHPObject{
ClassName: "User",
Properties: map[string]interface{}{
"id":       123,
"username": "john_doe",
},
}

serialized, _ := phpserialize.MarshalObject(phpObject)
// Output: O:4:"User":2:{s:2:"id";i:123;s:8:"username";s:8:"john_doe";}

// Unserialization
result, _ := phpserialize.Unmarshal(serialized)
phpObj := result.(phpserialize.PHPObject)
fmt.Printf("Class: %s, ID: %d\n", phpObj.ClassName, phpObj.Properties["id"])
// Output: Class: User, ID: 123
```

## API Reference and Options

The core functions are `Marshal` and `Unmarshal`. Both accept an optional list of Option interfaces for customization.

### Functions

| Function                                                          | Description                                       |
|-------------------------------------------------------------------|---------------------------------------------------|
| `Marshal(value interface{}, options ...Option) (string, error)-`  | Serializes a Go value to PHP format.              |
| `Unmarshal(data string, options ...Option) (interface{}, error)`  | Unserializes PHP data to Go values.               |
| `MarshalObject(obj PHPObject, options ...Option) (string, error)` | Dedicated function for serializing a `PHPObject`. |

### Helper Functions

| Function	                                                    | Description                                       |
|--------------------------------------------------------------|---------------------------------------------------|
| `IsValidMarshaled(data string) bool`                         | Checks if a string is valid PHP serialized data.  |
| `MustMarshal(value interface{}, options ...Option) string`   | Serializes and panics on error.                   |
| `MustUnmarshal(data string, options ...Option) interface{}`  | Unserializes and panics on error.                 |

## Options ⚙

Options are created using factory functions and implement the unified `Option` interface.

### `WithMaxDepth(depth int)`

Limits the nesting depth to prevent stack overflow or DoS attacks from deeply nested, malicious data.

| Context    | `depth=0` Behavior                       | Default |
|------------|------------------------------------------|---------|
| Marshal    | Unlimited (matches PHP's behavior)       | `0`     |
| Unmarshal  | Uses PHP's default maximum depth of 4096 | `4096`  |

```go
// Limit serialization/unserialization to 10 levels
data, err := phpserialize.Marshal(nestedData, phpserialize.WithMaxDepth(10))
result, err := phpserialize.Unmarshal(data, phpserialize.WithMaxDepth(10))
```

### `WithStrictPHP(strict bool)`

Enforces strict PHP compatibility rules (default: `true`). When set to `false`, it may allow Go-specific features
like `uint64` serialization (`u:<number>;`), which is not standard in PHP but can be useful for internal Go-to-Go data
transfer.

```go
// Allow Go-native uint64 serialization
data, _ := phpserialize.Marshal(uint64(math.MaxUint64), phpserialize.WithStrictPHP(false))
// Output: u:18446744073709551615; (Non-strict PHP format)
```

### `WithAllowedClasses(classes []string)`

**Security Feature**: Restricts which PHP classes can be un-serialized to prevent **POP chains** or other remote code
execution vulnerabilities from untrusted data.

1. Pass a slice of allowed class names (e.g., `[]string{"User", "Order"}`).
2. Pass `nil` or an empty slice (`[]string{}`) to disable all object unserialization entirely (
   like `PHP's allowed_classes = false`).

```go
// Only allow "User" and "Admin" classes
result, err := phpserialize.Unmarshal(data,
phpserialize.WithAllowedClasses([]string{"User", "Admin"}))

// Disable all object unserialization
result, err := phpserialize.Unmarshal(data,
phpserialize.WithAllowedClasses(nil))
```

## Type Mapping ↔️

### PHP to Go Type Conversion (Unmarshal)

| PHP Type              | Go Type                  | Notes                                                 |
|-----------------------|--------------------------|-------------------------------------------------------|
| `null`                | `nil`                    ||
| `boolean`             | `bool`                   ||
| `integer`             | `int64`                  | All integers become 64-bit to prevent overflow.       |
| `float`/`double`      | `float64`                | Includes NaN, Inf, -Inf.                              |
| `string`              | `string`                 | Binary safe.                                          |
| `array` (sequential)  | `[]interface{}`          | Only if keys are sequential integers starting from 0. |
| `array` (associative) | `map[string]interface{}` | Any other array key structure.                        |
| `object`              | `phpserialize.PHPObject` | Contains ClassName and Properties.                    |

### Go to PHP Type Conversion (Marshal)

| Go Type                  | PHP Type            | PHP Serialized Format               |
|--------------------------|---------------------|-------------------------------------|
| `nil`                    | `null`              | `N;`                                |
| `bool`                   | `boolean`           | `b:0;` or `b:1;`                    |
| `int`, `int64`           | `integer`           | `i:<number>;`                       |
| `float64`                | `double`            | `d:<number>;`                       |
| `string`                 | `string`            | `s:<length>:"<string>";`            |
| `[]interface{}`          | `array`             | `a:<count>:{...} (indexed keys)`    |
| `map[string]interface{}` | `associative array` | `a:<count>:{...} (string/int keys)` |
| `phpserialize.PHPObject` | `object`            | `O:<len>:"<class>":<count>:{...}`   |