# datagen

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/dream-horizon-org/datagen)](https://goreportcard.com/report/github.com/dream-horizon-org/datagen)
[![Join our Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/f92f4bWp)

datagen is a tool to generate coherent, synthetic data generation from models expressed in a simple, declarative DSL.

## Watch the Demo

<div align="center">

[![Watch the video](https://img.youtube.com/vi/ly0DfzTup28/maxresdefault.jpg)](https://www.youtube.com/watch?v=ly0DfzTup28)

</div>

Salient features:
* A **declarative DSL** for defining data models with Go-like syntax
* **High performance** through transpilation to native Go code
* **Multiple output formats** (CSV, JSON, XML, stdout)
* **Database integration** with direct loading to MySQL
* **Model relationships** via cross-references using `self.datagen`
* **Tag-based filtering** for selective data generation
* **Built-in functions** for common data items

## Install

See the [Installation Guide](http://dream-horizon-org.github.io/datagen/introduction/getting-started/#installation) for detailed installation instructions.

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
datagenc gen user.dg -f csv -o ./output
```

this will generate a `user.csv` file in `output` directory with 100 user records.

## More information

* See the [Documentation](https://dream-horizon-org.github.io/datagen/) for details

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License, see [LICENSE](LICENSE).
