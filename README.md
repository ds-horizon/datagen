# datagen

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ds-horizon/datagen)](https://goreportcard.com/report/github.com/ds-horizon/datagen)
[![Join our Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/cvMa8HrN)

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

There are various ways of installing datagen.

### Option 1: Install via Go

Check your `$PATH`, and choose a directory you would like to place the `datagenc` compiler in.

```bash
echo $PATH
/Users/username/go/bin:/opt/homebrew/bin:/opt/homebrew/sbin
```

Say, you wish to place the binary in `/opt/homebrew/bin`;

```bash
export GOBIN=/opt/homebrew/bin
go install github.com/ds-horizon/datagen/cmd/datagenc@latest
```

##### Verify installation

```bash
datagenc --help
```

### Option 2: Install from Source

##### Clone the repository

```bash
git clone github.com/ds-horizon/datagen
```

##### Build the compiler

```bash
make build-compiler
```

For permanent access on Mac/Unix, add the binary to your path, or add the current directory to your path:

```bash
echo 'export PATH=$PATH:$(pwd)' >> ~/.bashrc  # for bash
echo 'export PATH=$PATH:$(pwd)' >> ~/.zshrc   # for zsh
```

For permanent access on Windows, add to your shell profile:

```powershell
echo '$env:PATH += ";C:\path\to\datagen"' >> $PROFILE
```

Now, source the `rc` files or fire up a new terminal window for the changes to take effect.

##### Verify installation

```bash
datagenc --help
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
datagenc gen user.dg -f csv -o ./output
```

this will generate a `user.csv` file in `output` directory with 100 user records.

## More information

* See the [Documentation](https://ds-horizon.github.io/datagen/) for details

## Contributing

Refer to [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License, see [LICENSE](LICENSE).
