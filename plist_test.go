package plist

import (
	"reflect"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<array>
	<dict>
		<key>_SPCommandLineArguments</key>
		<array>
			<string>/usr/sbin/system_profiler</string>
			<string>-nospawn</string>
			<string>-xml</string>
			<string>SPPowerDataType</string>
			<string>-detailLevel</string>
			<string>full</string>
		</array>
	</dict>
</array>
</plist>
`
	v, err := Read(strings.NewReader(xml))
	if err != nil {
		t.Fatal(err)
	}
	a, ok := v.(Array)
	if !ok {
		t.Fatal("should be array")
	}
	if len(a) != 1 {
		t.Fatal("should have 1 child")
	}
	v = a[0]
	d, ok := v.(Dict)
	if !ok {
		t.Fatal("should have dict")
	}
	if len(d) != 1 {
		t.Fatal("should have 1 key")
	}
	key := ""
	for k, _ := range d {
		key = k
		break
	}
	if key != "_SPCommandLineArguments" {
		t.Fatal("key should be _SPCommandLineArguments")
	}
	v, ok = d[key]
	if !ok {
		t.Fatal("dict should have value")
	}
	a, ok = v.(Array)
	if !ok {
		t.Fatal("should be array")
	}
	if len(a) != 6 {
		t.Fatal("should have 6 children")
	}
	for i, s := range []string{"/usr/sbin/system_profiler", "-nospawn", "-xml", "SPPowerDataType", "-detailLevel", "full"} {
		if !reflect.DeepEqual(a[i], s) {
			t.Fatalf("a[%d] should be %v", i, s)
		}
	}
}
