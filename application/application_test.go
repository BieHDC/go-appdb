package application_test

import (
	. "biehdc.webapp.applister/application"
	"strings"
	"testing"
)

type fakecsv struct {
	testname      string
	shouldsucceed bool
	expectederror string
	fakeentry     string
}

var tests []fakecsv = []fakecsv{
	{"Minimal Success", true, "", "Name;Example Application\nVersion;4.20\nUseability;3\nRosVersion;0.4.13-i386\nTags;sneed\n"},
	{"Missing Name", false, "missing Name", "Version;4.20\nUseability;3\nRosVersion;0.4.13-i386\nTags;sneed\n"},
	{"Missing Version", false, "missing Version", "Name;Example Application\nUseability;3\nRosVersion;0.4.13-i386\nTags;sneed\n"},
	{"Missing Useability", false, "missing Useability", "Name;Example Application\nVersion;4.20\nRosVersion;0.4.13-i386\nTags;sneed\n"},
	{"Missing RosVersion", false, "missing RosVersion", "Name;Example Application\nVersion;4.20\nUseability;3\nTags;sneed\n"},
	{"Missing Tags", false, "missing Tags", "Name;Example Application\nVersion;4.20\nUseability;3\nRosVersion;0.4.13-i386\n"},
	{"Useability not a number", false, "error converting number: Nullptr reason: strconv.ParseUint: parsing \"Nullptr\": invalid syntax", "Name;Example Application\nVersion;4.20\nUseability;Nullptr\nRosVersion;0.4.13-i386\nTags;sneed\n"},
	{"Unknown field", false, "unknown field: ThisFieldDoesNotExist", "ThisFieldDoesNotExist;This Value Doesnt Matter\n"},
	{"Downloads entry too short", false, "Downloads entry too short", "Name;Example Application\nVersion;4.20\nUseability;3\nRosVersion;0.4.13-i386\nTags;sneed\nDownloads;missinglink\n"},
	{"Screenshot entry too short", false, "Screenshots entry too short", "Name;Example Application\nVersion;4.20\nUseability;3\nRosVersion;0.4.13-i386\nTags;sneed\nScreenshots;missinglink\n"},
}

func TestApplicationFromCSV(t *testing.T) {
	for _, tt := range tests {
		tt := tt //You need this
		t.Run(tt.testname, func(t *testing.T) {
			_, err := ApplicationFromCSV(strings.NewReader(tt.fakeentry))
			if err != nil {
				if err.Error() != tt.expectederror {
					t.Errorf("expected >%v<, got >%v<", tt.expectederror, err.Error())
				}
			}
			if !tt.shouldsucceed && err == nil {
				t.Error("expected fail, got success")
			}
		})
	}

}
