### Metadata

`metadata` provides configuration options that control model behavior and organization.

## Overview

The metadata section is defined within a model and supports configuration options:

```go
model example {
  metadata {
    count: 10
    tags: {
      "service": "user",
      "team": "backend"
    }
  }
  
  fields {
    // field definitions
  }
  
  gens {
    // generator implementations
  }
}
```

## `count`

Sets a default number of records that will be generated when no explicit count is provided.

### Basic Usage
```go
metadata {
  count: 5
}
```

### Command Behavior
- **Without `-n` flag**: Uses the model's `metadata.count` value
- **With `-n <number>` flag**: Overrides the default for ALL selected models

## `tags`

Allows you to label models with string key-value pairs, which can be used to filter the models to generate the data for.

### Basic Usage
```go
metadata {
  tags: {
    "service": "delivery",
    "team": "backend",
    "environment": "production"
  }
}
```

### Filtering Rules
- Models must match **ALL** specified tags to be selected
- Tag values must match exactly (case-sensitive)
- Use comma-separated values for multiple tag filters
