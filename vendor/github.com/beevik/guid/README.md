[![Build Status](https://travis-ci.org/beevik/guid.svg?branch=master)](https://travis-ci.org/beevik/guid)
[![GoDoc](https://godoc.org/github.com/beevik/guid?status.svg)](https://godoc.org/github.com/beevik/guid)

guid
====

The guid package implements a 16-byte guid/uuid type. It supports parsing
of guid strings, validation of guids, and random generation of guids
according to [RFC-4122](http://www.ietf.org/rfc/rfc4122.txt).

See http://godoc.org/github.com/beevik/guid for the godoc-formatted API
documentation.

### Example: Parsing a guid

```go
g, err := guid.ParseString("67a23ff3-20be-4420-9274-d16f2833d595")
```

### Example: Generating a random guid

```go
g := guid.New()
```

### Example: Validating a guid string

```go
s0 := "67a23ff3-20be-4420-9274-d16f2833d595"
s1 := "67a23ff3-20be-4420-9274"
fmt.Println("s0 a guid?  ", guid.IsGuid(s0))
fmt.Println("s0 a guid?  ", guid.IsGuid(s1))
```

Output:
```
s0 a guid?  true
s1 a guid?  false
```

### Example: Converting a guid to a string

```go
for i := 0; i < 4; i++ {
	g := guid.New()
	fmt.Println("guid: %s  GUID: %s\n", g.String(), g.StringUpper())
}
```

Output:
```
guid: 9a5bb29c-cdcd-4b1b-a039-b88d1271ab4c  GUID: 9A5BB29C-CDCD-4B1B-A039-B88D1271AB4C
guid: efeeee74-aea6-4fa9-8037-8d3a3f883d3f  GUID: EFEEEE74-AEA6-4FA9-8037-8D3A3F883D3F
guid: 773730b7-7b5d-4ef0-80ed-f52617d5b688  GUID: 773730B7-7B5D-4EF0-80ED-F52617D5B688
guid: 256f523f-e65e-4451-9381-1d561abbc645  GUID: 256F523F-E65E-4451-9381-1D561ABBC645
```
