# Fields Examples

## Overview

The `fields` section declares the individual elements that make up the schema, and their corresponding types. Each field has a name, type, and corresponding generator function that produces values for that field.

## Basic Syntax

```
fields {
  field_name() return_type
  field_with_params(param1 type1, param2 type2) return_type
}
```

## Examples

### 1. [Single Field Model](./single_field_model/)
- Shows basic, static value generation


### 2. [Multi-Field Model](./multi_field_model/)
- Shows various primitive types (int, string, bool, float32, time.Time), and use of `iter` for sequential ID generation

## Common Patterns

### Sequential IDs
```go
func id() {
  return iter + 1  // Creates 1, 2, 3, 4...
}
```

### Random Values
```go
func amount() {
  return FloatBetween(10.0, 1000.0)
}
```
