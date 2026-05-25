package detectors

import "testing"

func TestParseCountdownSeconds(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "minute second format", text: "1:23", want: 83},
		{name: "chinese format", text: "2分05", want: 125},
		{name: "single number", text: "45", want: 45},
		{name: "invalid", text: "马上开奖", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCountdownSeconds(tt.text); got != tt.want {
				t.Fatalf("parseCountdownSeconds(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}
