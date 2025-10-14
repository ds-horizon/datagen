# datagen

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ds-horizon/datagen)](https://goreportcard.com/report/github.com/ds-horizon/datagen)

datagen is a tool to generate coherent, synthetic data generation from models expressed in a simple, declarative DSL.

Salient features:
* A **declarative DSL** for defining data models with Go-like syntax
* **High performance** through transpilation to native Go code
* **Multiple output formats** (CSV, JSON, XML, stdout)
* **Database integration** with direct loading to MySQL
* **Model relationships** via cross-references using `self.datagen`
* **Tag-based filtering** for selective data generation
* **Built-in functions** for common data items

## Install

There are various ways of installing datagen.

### Precompiled binaries

Precompiled binaries for released versions are available in the [releases](https://github.com/ds-horizon/datagen/releases) section. Using the latest production release binary is the recommended way of installing datagen.

### Building from source

To build DataGen from source code, you need Go 1.24.0 or later.

Start by cloning the repository:

```bash
git clone https://github.com/ds-horizon/datagen.git
cd datagen
```

You can use `make` to build the `datagen` binary:

```bash
make build-compiler
./datagen --help
```

## Usage

You can launch datagen for trying it out with:

```bash
# Create a simple model file
cat > user.dg << 'EOF'
model user {
  metadata { count: 100 }
  fields {
    id() int
    name() string
  }
  gens {
    func id() { return iter + 1 }
    func name() { return Name() }
  }
}
EOF

# Generate data
./datagen gen user.dg -f csv -o ./output
```

datagen will now generate a `user.csv` file in `output` directory with 100 user records.

## More information

* See the [Documentation](https://github.com/ds-horizon/datagen/discussions) for details

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License, see [LICENSE](LICENSE).
