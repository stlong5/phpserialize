package main

import (
	"encoding/json"
	"fmt"
	"github.com/stlong5/phpserialize"
	"math"
)

func main() {
	fmt.Println("=== PHP Serialize/Unserialize Examples ===")

	// Example 1: Basic Types
	example1()

	// Example 2: Strings with Special Characters
	example2()

	// Example 3: Arrays and Maps
	example3()

	// Example 4: PHP Objects
	example4()

	// Example 5: Complex Nested Structures
	example5()

	// Example 6: Depth Limiting
	example6()

	// Example 7: Security Features
	example7()

	// Example 8: Special Float Values
	example8()

	// Example 9: Database Integration
	example9()

	// Example 10: JSON within Serialized Data
	example10()

	// Example 11: Migration Helper
	example11()
}

func example1() {
	fmt.Println("Example 1: Basic Types")
	fmt.Println("----------------------")

	// Null
	data, _ := phpserialize.Marshal(nil)
	fmt.Printf("nil -> %s\n", data)

	// Boolean
	data, _ = phpserialize.Marshal(true)
	fmt.Printf("true -> %s\n", data)
	data, _ = phpserialize.Marshal(false)
	fmt.Printf("false -> %s\n", data)

	// Integer
	data, _ = phpserialize.Marshal(42)
	fmt.Printf("42 -> %s\n", data)

	// Float
	data, _ = phpserialize.Marshal(3.14159)
	fmt.Printf("3.14159 -> %s\n", data)

	// String
	data, _ = phpserialize.Marshal("Hello, PHP!")
	fmt.Printf("\"Hello, PHP!\" -> %s\n", data)

	// Unserialize back
	result, _ := phpserialize.Unmarshal(data)
	fmt.Printf("Unserialized: %v (type: %T)\n", result, result)

	fmt.Println()
}

func example2() {
	fmt.Println("Example 2: Strings with Special Characters")
	fmt.Println("-------------------------------------------")

	// String with quotes
	str1 := `She said, "It's amazing!"`
	serialized1, _ := phpserialize.Marshal(str1)
	fmt.Printf("Original: %s\n", str1)
	fmt.Printf("Serialized: %s\n", serialized1)

	// String with newlines and tabs
	str2 := "Line 1\nLine 2\tTabbed"
	serialized2, _ := phpserialize.Marshal(str2)
	fmt.Printf("\nOriginal with newlines: %q\n", str2)
	fmt.Printf("Serialized: %s\n", serialized2)

	// Unicode string
	str3 := "Hello 世界 🌍"
	serialized3, _ := phpserialize.Marshal(str3)
	fmt.Printf("\nUnicode: %s\n", str3)
	fmt.Printf("Serialized: %s\n", serialized3)
	fmt.Printf("Byte length: %d (not character count)\n", phpserialize.GetStringLength(str3))

	// Verify round-trip
	unserialized, _ := phpserialize.Unmarshal(serialized3)
	fmt.Printf("Round-trip OK: %v\n", unserialized == str3)

	fmt.Println()
}

func example3() {
	fmt.Println("Example 3: Arrays and Maps")
	fmt.Println("---------------------------")

	// Indexed array (becomes slice in Go)
	arr := []interface{}{"apple", "banana", "cherry"}
	data, _ := phpserialize.Marshal(arr)
	fmt.Printf("Array: %s\n", data)

	result, _ := phpserialize.Unmarshal(data)
	fmt.Printf("Unserialized as slice: %v (type: %T)\n", result, result)

	// Associative array (map in Go)
	m := map[string]interface{}{
		"name":   "John",
		"age":    30,
		"city":   "New York",
		"active": true,
	}
	data, _ = phpserialize.Marshal(m)
	fmt.Printf("\nMap: %s\n", data)

	result, _ = phpserialize.Unmarshal(data)
	if resultMap, ok := result.(map[string]interface{}); ok {
		fmt.Printf("Name: %v, Age: %v\n", resultMap["name"], resultMap["age"])
	}

	fmt.Println()
}

func example4() {
	fmt.Println("Example 4: PHP Objects")
	fmt.Println("----------------------")

	// Create a PHP object
	obj := phpserialize.PHPObject{
		ClassName: "User",
		Properties: map[string]interface{}{
			"id":         123,
			"username":   "john_doe",
			"email":      "john@example.com",
			"created_at": "2024-01-15 08:30:00",
			"is_active":  true,
			"settings": map[string]interface{}{
				"newsletter": true,
				"privacy":    "public",
			},
		},
	}

	data, _ := phpserialize.MarshalObject(obj)
	fmt.Printf("Serialized: %s\n", data)

	// Unserialize
	result, _ := phpserialize.Unmarshal(data)
	if phpObj, ok := result.(phpserialize.PHPObject); ok {
		fmt.Printf("\nClass: %s\n", phpObj.ClassName)
		fmt.Printf("Username: %v\n", phpObj.Properties["username"])
		fmt.Printf("Settings: %v\n", phpObj.Properties["settings"])
	}

	fmt.Println()
}

func example5() {
	fmt.Println("Example 5: Complex Nested Structures")
	fmt.Println("-------------------------------------")

	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    1001,
			"name":  "Alice Smith",
			"roles": []interface{}{"admin", "editor", "moderator"},
			"metadata": map[string]interface{}{
				"last_login": "2025-01-15 10:30:00",
				"preferences": map[string]interface{}{
					"theme":         "dark",
					"notifications": true,
				},
				"stats": map[string]interface{}{
					"posts":    156,
					"comments": 892,
				},
			},
		},
		"timestamp": 1736899200,
		"version":   "2.0.1",
	}

	serialized, _ := phpserialize.Marshal(complexData)
	fmt.Printf("Serialized (%d bytes)\n", len(serialized))
	fmt.Printf("First 150 chars: %.150s...\n", serialized)

	// Unserialize and extract data
	unserialized, _ := phpserialize.Unmarshal(serialized)
	if data, ok := unserialized.(map[string]interface{}); ok {
		if user, ok := data["user"].(map[string]interface{}); ok {
			fmt.Printf("\nUser name: %v\n", user["name"])
			if metadata, ok := user["metadata"].(map[string]interface{}); ok {
				if prefs, ok := metadata["preferences"].(map[string]interface{}); ok {
					fmt.Printf("Theme: %v\n", prefs["theme"])
				}
			}
		}
	}

	fmt.Println()
}

func example6() {
	fmt.Println("Example 6: Depth Limiting")
	fmt.Println("-------------------------")

	// Create nested data
	data := []interface{}{
		1,
		[]interface{}{2, 3},
		[]interface{}{4, []interface{}{5, 6}},
	}

	// maxDepth=5: levels 0,1,2,3,4 allowed
	serialized, err := phpserialize.Marshal(data, phpserialize.WithMaxDepth(5))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Serialized: %s\n", serialized)

		// Same limit for unmarshal
		result, err := phpserialize.Unmarshal(serialized, phpserialize.WithMaxDepth(5))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Unserialized: %v\n", result)
		}
	}

	// Test depth limit exceeded
	// maxDepth=2 means levels 0,1 allowed, level 2 blocked (>= 2)
	deepData := []interface{}{
		[]interface{}{
			[]interface{}{
				"level 2 - blocked",
			},
		},
	}

	_, err = phpserialize.Marshal(deepData, phpserialize.WithMaxDepth(2))
	if err != nil {
		fmt.Printf("\nExpected error (depth >= 2): %v\n", err)
	}

	// Test unmarshal depth limit
	deepSerialized := `a:1:{i:0;a:1:{i:0;a:1:{i:0;s:5:"deep";}}}`
	_, err = phpserialize.Unmarshal(deepSerialized, phpserialize.WithMaxDepth(2))
	if err != nil {
		fmt.Printf("Expected error on unmarshal: %v\n", err)
	}

	fmt.Println()
}

func example7() {
	fmt.Println("Example 7: Security Features")
	fmt.Println("-----------------------------")

	// Class filtering
	objectData := `O:4:"User":2:{s:2:"id";i:123;s:4:"name";s:4:"John";}`

	// Allow all classes (default)
	result, _ := phpserialize.Unmarshal(objectData)
	fmt.Printf("Allow all: %v\n", result)

	// Only allow specific classes
	result, err := phpserialize.Unmarshal(objectData,
		phpserialize.WithAllowedClasses([]string{"User", "Admin"}))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Allowed User: %v\n", result)
	}

	// Disallow all classes
	_, err = phpserialize.Unmarshal(objectData,
		phpserialize.WithAllowedClasses(nil))
	if err != nil {
		fmt.Printf("Expected error (no classes): %v\n", err)
	}

	// Combine security options
	result, _ = phpserialize.Unmarshal(objectData,
		phpserialize.WithMaxDepth(100),
		phpserialize.WithAllowedClasses([]string{"User"}))
	fmt.Printf("\nWith combined security: %v\n", result)

	fmt.Println()
}

func example8() {
	fmt.Println("Example 8: Special Float Values")
	fmt.Println("--------------------------------")

	// NaN
	data, _ := phpserialize.Marshal(math.NaN())
	fmt.Printf("NaN: %s\n", data)
	result, _ := phpserialize.Unmarshal(data)
	fmt.Printf("Unserialized: IsNaN=%v\n", math.IsNaN(result.(float64)))

	// Infinity
	data, _ = phpserialize.Marshal(math.Inf(1))
	fmt.Printf("INF: %s\n", data)

	data, _ = phpserialize.Marshal(math.Inf(-1))
	fmt.Printf("-INF: %s\n", data)

	// Large number (decimal format, not scientific)
	largeNum := 123456789.123456789
	data, _ = phpserialize.Marshal(largeNum)
	fmt.Printf("\nLarge number: %s (decimal, not 1.234e+08)\n", data)

	fmt.Println()
}

func example9() {
	fmt.Println("Example 9: Database Integration (Simulated)")
	fmt.Println("-------------------------------------------")

	// Simulate PHP serialized data from database
	phpFromDB := `a:4:{s:10:"site_title";s:10:"My Website";s:7:"version";s:5:"1.2.3";s:8:"settings";a:3:{s:5:"theme";s:4:"dark";s:8:"language";s:2:"en";s:8:"timezone";s:3:"UTC";}s:5:"users";i:1250;}`

	fmt.Printf("From DB: %.80s...\n\n", phpFromDB)

	// Unserialize
	config, _ := phpserialize.Unmarshal(phpFromDB)
	fmt.Printf("Config: %v\n", config)

	// Modify
	if configMap, ok := config.(map[string]interface{}); ok {
		configMap["last_updated"] = "2025-01-15"
		configMap["users"] = int64(1251)

		if settings, ok := configMap["settings"].(map[string]interface{}); ok {
			settings["theme"] = "light"
		}

		// Serialize back for database
		newData, _ := phpserialize.Marshal(configMap)
		fmt.Printf("\nUpdated for DB: %.100s...\n", newData)
		fmt.Println("SQL: UPDATE settings SET value = ? WHERE key = 'config'")
	}

	fmt.Println()
}

func example10() {
	fmt.Println("Example 10: JSON within Serialized Data")
	fmt.Println("----------------------------------------")

	// Common scenario: PHP stores JSON, then serializes it
	jsonData := map[string]interface{}{
		"name":    "Product ABC",
		"price":   99.99,
		"tags":    []interface{}{"electronics", "gadgets"},
		"inStock": true,
	}

	// Convert to JSON (as PHP might do)
	jsonBytes, _ := json.Marshal(jsonData)
	jsonString := string(jsonBytes)
	fmt.Printf("JSON: %s\n", jsonString)

	// Serialize the JSON string
	serialized, _ := phpserialize.Marshal(jsonString)
	fmt.Printf("Serialized: %s\n", serialized)

	// Complex structure with JSON inside
	complexData := map[string]interface{}{
		"id":          1,
		"type":        "product",
		"json_data":   jsonString,
		"description": `Product with "quotes" and 'apostrophes'`,
	}

	serializedComplex, _ := phpserialize.Marshal(complexData)
	fmt.Printf("\nComplex: %.120s...\n", serializedComplex)

	// Extract and parse JSON
	unserialized, _ := phpserialize.Unmarshal(serializedComplex)
	if data, ok := unserialized.(map[string]interface{}); ok {
		if jsonStr, ok := data["json_data"].(string); ok {
			var parsed map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &parsed)
			fmt.Printf("\nExtracted JSON: %v\n", parsed)
			fmt.Printf("Product name: %v\n", parsed["name"])
		}
	}

	fmt.Println()
}

func example11() {
	fmt.Println("Example 11: Migration Helper")
	fmt.Println("-----------------------------")

	// Helper to migrate PHP session to Go format
	migrateSession := func(phpSession string) (map[string]interface{}, error) {
		data, err := phpserialize.Unmarshal(phpSession)
		if err != nil {
			return nil, err
		}

		if sessionMap, ok := data.(map[string]interface{}); ok {
			// Convert PHP keys to Go conventions
			goSession := make(map[string]interface{})
			if val, ok := sessionMap["user_id"]; ok {
				goSession["userID"] = val
			}
			if val, ok := sessionMap["username"]; ok {
				goSession["username"] = val
			}
			if val, ok := sessionMap["last_activity"]; ok {
				goSession["lastActivity"] = val
			}
			return goSession, nil
		}

		return nil, fmt.Errorf("invalid format")
	}

	// Example PHP session
	phpSession := `a:3:{s:7:"user_id";i:42;s:8:"username";s:8:"john_doe";s:13:"last_activity";s:19:"2025-01-15 14:30:00";}`

	fmt.Printf("PHP session: %s\n", phpSession)

	goSession, _ := migrateSession(phpSession)
	fmt.Printf("Go session: %v\n", goSession)

	// Batch processing
	fmt.Println("\nBatch processing:")
	batch := []string{
		`s:5:"hello";`,
		`i:42;`,
		`a:2:{i:0;s:3:"foo";i:1;s:3:"bar";}`,
		`b:1;`,
		`N;`,
	}

	for i, item := range batch {
		if phpserialize.IsValidMarshaled(item) {
			result, _ := phpserialize.Unmarshal(item)
			fmt.Printf("  %d: %v (type: %T)\n", i+1, result, result)
		} else {
			fmt.Printf("  %d: Invalid\n", i+1)
		}
	}

	fmt.Println()
}
