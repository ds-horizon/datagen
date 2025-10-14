# Datagen DSL Specification

## Introduction

Datagen is a declarative DSL (domain specific language) that allows the user to encode the schema of the entity they wish to generate. It is composed of separate sections, that allow the user to specify different aspects of the schema, which are detailed in this document. The datagen compiler ingests this datagen code (in `.dg` files) and transpiles it into go code.

Salient features:
- Declarative model definitions
- Custom generator functions
- Metadata-driven configuration
- Cross-model references
- Integration with Go's type system

## Lexical Elements

### Keywords

The following are reserved keywords in the Datagen DSL:

- `model` - Declares a data model
- `fields` - Defines field specifications
- `gens` - Defines generator functions
- `calls` - Defines function calls for field initialization
- `misc` - Injects arbitrary Go code
- `func` - Declares a function
- `metadata` - Defines model metadata
    - `count` - Specifies default record count
    - `tags` - Defines model tags

*Unless otherwise specified, the definitions of identifiers, literals, punctuators are is the same as in the [Go spec](https://go.dev/ref/spec)*

## Syntax

### Grammar Overview

The Datagen DSL follows this top-level grammar:

```
%start main
main: "model" identifier "{" main_body "}"
main_body: fields_section
           | misc_section
           | metadata_section
           | calls_section
           | gen_fns_section

```

## `model`

A model declaration is the top level structure for encoding the schema. It consists of a model name followed by a block containing various sections.
Each `.dg` file contains a single model declaration, with the file name matching the model name.

### Syntax

```
main: "model" identifier "{" main_body "}"
```

### Example

```datagen
model user_profile {}
```

Here, we define an entity called `user_profile` which has no fields (not a very useful entity).

## `metadata`

The metadata section provides configuration information for the model, including default record counts and tags for filtering.
`count` is an integer, and `tags` is a set of string key-value pairs.

### Syntax

```
metadata_section: "metadata" "{" metadata_body "}"
metadata_body: count_entry metadata_body
               | tags_entry metadata_body
               | // empty
count_entry: "count" ":" COUNT_INT
tags_entry: "tags" ":" "{" tags_body "}"
tags_body: "<key>" ":" <value> "," tags_body
           | // empty
```

### Example

```datagen
metadata {
    count: 100
    tags: {
        "service": "user",
        "team": "backend"
    }
}
```

## `fields`

The fields section declares the individual elements that make up the schema, and their corresponding types.

### Syntax

```
fields_section: "fields" "{" fields_body "}"
fields_body: field_declaration fields_body 
             | // empty

field_declaration: Identifier "(" [GoType]")" GoType
```

### Example

```datagen
fields {
    id() int
    name() string
    email() string
    is_active() bool
    rating() float32
    created_at() time.Time
    metadata() map[string]string
}
```

### Supported Types

The Datagen DSL supports all valid Go types, including:

- **Primitive types**: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `string`, `bool`, `float32`, `float64`, `complex64`, `complex128`
- **Built-in types**: `time.Time`, `time.Duration`
- **Composite types**: `map`, `slice`, `array`, `struct`, `interface{}`
- **Pointer types**: `*Type`
- **Custom types**: Any Go type defined in the `misc` section or imported from other Go packages


## `gens`

The generator functions section is composed of a series of generator functions. Each generator function corresponds to a field and returns a value of the appropriate type as required for that field.

### Syntax

```
gen_fns_section: "gens" "{" gen_fns "}"
gen_fns: "func" FuncNameLiteral "(" [parameter_list] ")" "{" <GoCode> "}" gen_fns 
         | //empty
parameter_list = parameter { "," parameter }
parameter = Identifier Type
```
### Example

```datagen
gens {
    func age() {
      return 42
    }
	
    func set_name() {
      return "user_data"
    }	

    func address() {
       d := &Address{
          City: Sentence(10),
          Zip: Sentence(10),
          Desc: Sentence(10),
        }
      return ToJSON(d)
    }
}
```

## `calls`

The calls section provides initialization arguments to the generator functions. These need to be specified only for those field generator functions that accept parameters.

### Syntax

```
calls_section = "calls" "{" call_list "}"

call_list = call call_list

call = field_name "(" [argument_list] ")"
field_name = Identifier
argument_list = <GoExpression> { "," <GoExpression> }
```

### Example

```datagen
calls {
    created_at("2023-01-01 00:00:00", "2023-12-31 23:59:59")
}
```

## `misc`

The misc section allows injection of arbitrary Go code into the generated package. This is useful for defining custom types, constants, and helper functions.

### Syntax

```
MiscSection = "misc" "{" <GoCode> "}"
```

### Example

```datagen
misc {
    type Address struct {
        Line1 string `json:"line1"`
        City string `json:"city"`
    }
}
```
