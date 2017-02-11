package plist

import (
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
	if len(a) == 0 {
		t.Fatal("should not be empty")
	}
}
