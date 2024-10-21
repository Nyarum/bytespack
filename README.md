## Code Generation for Encode/Decode Methods of binary protocol

Generate encode and decode methods for your structs to avoid using reflection with custom binary protocols.

### How to Use

1. Install the code generation tool:
   ```sh
   go install github.com/Nyarum/bytespack/cmd/diho_bytes_generate
   ```

2. Add a `go:generate` comment to your struct definition:
   ```go
   //go:generate diho_bytes_generate packet.go
    type Header struct {
        Pass uint8
    }

    type InternalStruct struct {
        Test uint16
    }

    //go:generate diho_bytes_generate packet.go
    type Packet struct {
        Header          `dbg:"ignore"`
        ID              uint16
        Name            string
        Level           uint32
        Health          uint8
        Mana            uint16
        Bro             int8
        Bro2            int16
        Bro3            int32
        Bro4            int64
        BytesField      []byte
        InternalStruct  InternalStruct
        ArrayTest       [5]uint16
        SliceTest       []uint16
        InternalStructs [2]InternalStruct
    }

    func (p *Packet) Filter(ctx context.Context) bool {
        return false
    }
   ```

3. Run the `go generate` command in your terminal:
   ```sh
   go generate ./...
   ```

This will produce the generated code in the same folder as your struct.

### Example of Generated Code

The generated code (for decode one) will look something like this:

```go
func (p *Header) Decode(ctx context.Context, buf []byte, endian binary.ByteOrder) error {
	var err error
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, endian, &p.Pass)
	if err != nil {
		return err
	}
	return nil
}
func (p *InternalStruct) Decode(ctx context.Context, buf []byte, endian binary.ByteOrder) error {
	var err error
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, endian, &p.Test)
	if err != nil {
		return err
	}
	return nil
}
func (p *Packet) Decode(ctx context.Context, buf []byte, endian binary.ByteOrder) error {
	var err error
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, endian, &p.ID)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	p.Name, err = utils.ReadStringNull(reader)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Level)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Health)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Mana)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Bro)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Bro2)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Bro3)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	err = binary.Read(reader, endian, &p.Bro4)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	p.BytesField, err = utils.ReadBytes(reader)
	if err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	if err = (&p.InternalStruct).Decode(ctx, buf, endian); err != nil {
		return err
	}
	if p.Filter(ctx) == true {
		return err
	}
	for k := range p.ArrayTest {
		var tempValue uint16
		if err = binary.Read(reader, endian, &tempValue); err != nil {
			return err
		}
		p.ArrayTest[k] = tempValue
	}
	if p.Filter(ctx) == true {
		return err
	}
	for k := range p.SliceTest {
		var tempValue uint16
		if err = binary.Read(reader, endian, &tempValue); err != nil {
			return err
		}
		p.SliceTest[k] = tempValue
	}
	if p.Filter(ctx) == true {
		return err
	}
	for k := range p.InternalStructs {
		if err = (&p.InternalStructs[k]).Decode(ctx, buf, endian); err != nil {
			return err
		}
	}
	if p.Filter(ctx) == true {
		return err
	}
	return nil
}
```

The generated code (for encode one) will look something like this:

```go
func (p *Header) Encode(ctx context.Context, endian binary.ByteOrder) ([]byte, error) {
	newBuf := bytebufferpool.Get()
	defer bytebufferpool.Put(newBuf)
	var err error
	err = binary.Write(newBuf, endian, p.Pass)
	if err != nil {
		return nil, err
	}
	return utils.Clone(newBuf), nil
}
func (p *InternalStruct) Encode(ctx context.Context, endian binary.ByteOrder) ([]byte, error) {
	newBuf := bytebufferpool.Get()
	defer bytebufferpool.Put(newBuf)
	var err error
	err = binary.Write(newBuf, endian, p.Test)
	if err != nil {
		return nil, err
	}
	return utils.Clone(newBuf), nil
}
func (p *Packet) Encode(ctx context.Context, endian binary.ByteOrder) ([]byte, error) {
	newBuf := bytebufferpool.Get()
	defer bytebufferpool.Put(newBuf)
	var err error
	err = binary.Write(newBuf, endian, p.ID)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = utils.WriteStringNull(newBuf, p.Name)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Level)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Health)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Mana)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Bro)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Bro2)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Bro3)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = binary.Write(newBuf, endian, p.Bro4)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	err = utils.WriteBytes(newBuf, p.BytesField)
	if err != nil {
		return nil, err
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	if encodeBuf, err := p.InternalStruct.Encode(ctx, endian); err != nil {
		return nil, err
	} else {
		newBuf.Write(encodeBuf)
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	for _, v := range p.ArrayTest {
		if err = binary.Write(newBuf, endian, v); err != nil {
			return nil, err
		}
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	for _, v := range p.SliceTest {
		if err = binary.Write(newBuf, endian, v); err != nil {
			return nil, err
		}
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	for _, v := range p.InternalStructs {
		if encodeBuf, err := v.Encode(ctx, endian); err != nil {
			return nil, err
		} else {
			newBuf.Write(encodeBuf)
		}
	}
	if p.Filter(ctx) == true {
		return utils.Clone(newBuf), nil
	}
	return utils.Clone(newBuf), nil
}
```
