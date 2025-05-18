# BytesPack

[![Go Version](https://img.shields.io/github/go-mod/go-version/Nyarum/bytespack)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

BytesPack is a powerful Go code generation tool that automatically generates encode/decode methods for your structs to handle binary protocols without using reflection. This makes binary protocol handling both faster and safer.

## Features

- 🚀 Zero reflection for better performance
- 🛠 Automatic code generation for encode/decode methods
- 💪 Support for various Go types including:
  - Basic types (uint8/16/32/64, int8/16/32/64)
  - Strings with null termination
  - Byte slices
  - Arrays and slices
  - Nested structs
- 🔧 Custom endianness support via struct tags
- 🎯 Field filtering capabilities
- 📦 Uses efficient byte buffer pool for better memory management

## Installation

```bash
go install github.com/Nyarum/bytespack/cmd/diho_bytes_generate@latest
```

## Quick Start

1. Define your struct with the `go:generate` directive:

```go
//go:generate diho_bytes_generate packet.go
type Packet struct {
    ID     uint16
    Name   string
    Level  uint32
    Health uint8
}
```

2. Run the code generation:

```bash
go generate ./...
```

This will create two files:
- `packet_encode.gen.go`: Contains the encoding logic
- `packet_decode.gen.go`: Contains the decoding logic

## Advanced Usage

### Struct Tags

BytesPack supports custom behavior through struct tags:

- `dbg:"ignore"` - Skip this field during encoding/decoding
- `dbg:"little"` - Use little-endian encoding for this field
- `dbg:"fieldName==value"` - Conditional encoding/decoding based on other field values

Example:

```go
type Packet struct {
    Header        `dbg:"ignore,little"`
    ID            uint16
    OptionalField uint32 `dbg:"ID==1"` // Only encoded/decoded if ID equals 1
}
```

### Custom Filtering

You can implement custom filtering logic by adding a Filter method to your struct:

```go
func (p *Packet) Filter(ctx context.Context, fieldName string) bool {
    // Return true to skip encoding/decoding of the current field
    return false
}
```

## Project Structure

```
bytespack/
├── cmd/
│   └── diho_bytes_generate/    # Code generation tool
├── customtypes/               # Custom type definitions
├── example/                   # Usage examples
├── generate/                  # Code generation logic
├── parse/                     # Struct parsing logic
└── utils/                     # Utility functions
```

## Performance

BytesPack generates code that:
- Avoids reflection completely
- Uses efficient byte buffer pooling
- Minimizes allocations
- Provides predictable performance

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
