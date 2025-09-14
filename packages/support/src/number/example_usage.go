// This file demonstrates the new error handling approach in the number package
// after decoupling from hashing exceptions.

package number

import (
	"fmt"
	"log"
	"math"
)

// ExampleUsage demonstrates how the conversion functions now return
// descriptive fmt errors instead of generic exception errors.
func ExampleUsage() {
	fmt.Println("=== Number Package - Decoupled Error Handling Example ===")
	
	// Example 1: Valid conversions
	fmt.Println("\n--- Valid Conversions ---")
	
	if val, err := ToUint32(42); err == nil {
		fmt.Printf("ToUint32(42) = %d ✓\n", val)
	}
	
	if val, err := ToUint8(255); err == nil {
		fmt.Printf("ToUint8(255) = %d ✓\n", val)
	}
	
	if val, err := ToInt(-42); err == nil {
		fmt.Printf("ToInt(-42) = %d ✓\n", val)
	}
	
	// Example 2: Invalid conversions with descriptive errors
	fmt.Println("\n--- Invalid Conversions (showing descriptive errors) ---")
	
	// Out of range for uint32
	if _, err := ToUint32(-1); err != nil {
		fmt.Printf("ToUint32(-1) error: %v\n", err)
	}
	
	// Out of range for uint8
	if _, err := ToUint8(256); err != nil {
		fmt.Printf("ToUint8(256) error: %v\n", err)
	}
	
	// Out of range for int (using a very large uint32)
	if _, err := ToInt(uint32(math.MaxUint32)); err != nil {
		fmt.Printf("ToInt(MaxUint32) error: %v\n", err)
	}
	
	// Unsupported type
	if _, err := ToUint32("not a number"); err != nil {
		fmt.Printf("ToUint32(string) error: %v\n", err)
	}
	
	// Non-integer float
	if _, err := ToUint8(3.14); err != nil {
		fmt.Printf("ToUint8(3.14) error: %v\n", err)
	}
	
	// Example 3: Error handling pattern
	fmt.Println("\n--- Recommended Error Handling Pattern ---")
	
	value := interface{}(1000)
	if result, err := ToUint8(value); err != nil {
		// Before: would get generic "ErrInvalidOptions"
		// Now: get descriptive error messages
		log.Printf("Conversion failed with details: %v", err)
		
		// You can still handle different error cases if needed
		switch {
		case err.Error() == "unsupported type string for uint8 conversion":
			fmt.Println("Handle unsupported type case")
		default:
			fmt.Printf("Handle other conversion error: %v\n", err)
		}
	} else {
		fmt.Printf("Conversion successful: %d\n", result)
	}
	
	fmt.Println("\n=== Benefits of Decoupling ===")
	fmt.Println("✓ No dependency on hashing/exceptions package")
	fmt.Println("✓ Descriptive error messages with context")
	fmt.Println("✓ Standard Go error handling patterns")
	fmt.Println("✓ Better debugging and error logging")
	fmt.Println("✓ More maintainable and testable code")
}

// DemonstrateErrorTypes shows the different types of errors that can occur
func DemonstrateErrorTypes() {
	fmt.Println("\n=== Error Type Examples ===")
	
	testCases := []struct {
		name     string
		function string
		value    interface{}
	}{
		{"Negative int to uint32", "ToUint32", -5},
		{"Large int64 to uint32", "ToUint32", int64(math.MaxUint32) + 1},
		{"Float with decimals to uint8", "ToUint8", 3.7},
		{"String to int", "ToInt", "hello"},
		{"Large uint32 to int", "ToInt", uint32(math.MaxUint32)},
	}
	
	for _, tc := range testCases {
		fmt.Printf("\nTest: %s\n", tc.name)
		
		switch tc.function {
		case "ToUint32":
			if _, err := ToUint32(tc.value); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		case "ToUint8":
			if _, err := ToUint8(tc.value); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		case "ToInt":
			if _, err := ToInt(tc.value); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		}
	}
}

// CompareOldVsNew shows the difference between old and new error handling
func CompareOldVsNew() {
	fmt.Println("\n=== Before vs After Comparison ===")
	
	fmt.Println("BEFORE (with hashing exceptions):")
	fmt.Println("  import \"govel/hashing/exceptions\"")
	fmt.Println("  return 0, exceptions.ErrInvalidOptions  // Generic error")
	fmt.Println("  // Error message: just \"invalid options\"")
	
	fmt.Println("\nAFTER (with fmt errors):")
	fmt.Println("  import \"fmt\"")
	fmt.Println("  return 0, fmt.Errorf(\"value %d is out of range for uint32 (0-%d)\", v, math.MaxUint32)")
	fmt.Println("  // Error message: \"value -5 is out of range for uint32 (0-4294967295)\"")
	
	fmt.Println("\nBenefits:")
	fmt.Println("  ✓ Self-contained package (no external dependencies)")
	fmt.Println("  ✓ Descriptive error messages with actual values")
	fmt.Println("  ✓ Standard Go idioms")
	fmt.Println("  ✓ Better debugging experience")
}