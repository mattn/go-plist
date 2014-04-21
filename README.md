# go-plist

plist parser

## Usage

```go
v, err := plist.Read(f)
if err != nil {
	log.Fatal(err)
}

tree := v.(plist.Dict)
for _, t := range tree["Tracks"].(plist.Dict) {
	if item, ok := t.(plist.Dict); ok {
		fmt.Println(item["Name"])
	}
}
```

## Installation

```
go get github.com/mattn/go-plist
```

# License

MIT

# Author

Yasuhiro Matsumoto (a.k.a mattn)
