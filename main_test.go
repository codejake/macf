package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNormalizeMAC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "colon-delimited", input: "c6:89:f2:d2:DC:3E", want: "c689f2d2dc3e"},
		{name: "hyphen-delimited", input: "c6-89-f2-d2-dc-3e", want: "c689f2d2dc3e"},
		{name: "cisco-style", input: "c689.f2d2.dc3e", want: "c689f2d2dc3e"},
		{name: "plain-hex", input: "C689F2D2DC3E", want: "c689f2d2dc3e"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeMAC(tc.input)
			if err != nil {
				t.Fatalf("normalizeMAC(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("normalizeMAC(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestNormalizeMACRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	invalid := []string{
		"",
		"c6:89:f2:d2:dc:3",
		"c6:89:f2:d2:dc:3e:4d:ab",
		"hello-c689-f2d2-dc3e",
	}

	for _, input := range invalid {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			if _, err := normalizeMAC(input); err == nil {
				t.Fatalf("normalizeMAC(%q) unexpectedly succeeded", input)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    formatStyle
		wantErr bool
	}{
		{input: "cisco", want: formatCisco},
		{input: "colon", want: formatColon},
		{input: "hyphen", want: formatHyphen},
		{input: "none", want: formatNone},
		{input: "bogus", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got, err := parseFormat(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseFormat(%q) unexpectedly succeeded", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFormat(%q) returned error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("parseFormat(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		stdin      string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:     "all formats",
			args:     []string{"c6:89:f2:d2:DC:3E"},
			wantCode: 0,
			wantStdout: strings.Join([]string{
				"c689.f2d2.dc3e",
				"C689F2D2DC3E",
				"c6:89:f2:d2:dc:3e",
				"C6:89:F2:D2:DC:3E",
				"C6-89-F2-D2-DC-3E",
				"",
			}, "\n"),
		},
		{
			name:       "stdin fallback",
			stdin:      "c6:89:f2:d2:DC:3E\n",
			wantCode:   0,
			wantStdout: "c689.f2d2.dc3e\nC689F2D2DC3E\nc6:89:f2:d2:dc:3e\nC6:89:F2:D2:DC:3E\nC6-89-F2-D2-DC-3E\n",
		},
		{
			name:       "none uppercase",
			args:       []string{"--format", "none", "--uppercase", "c6:89:f2:d2:DC:3E"},
			wantCode:   0,
			wantStdout: "C689F2D2DC3E\n",
		},
		{
			name:       "invalid format",
			args:       []string{"--format", "bogus", "c6:89:f2:d2:DC:3E"},
			wantCode:   1,
			wantStderr: "Error: invalid format \"bogus\" (must be one of: cisco, colon, hyphen, none)\n",
		},
		{
			name:       "invalid mac",
			args:       []string{"c6:89:f2:d2:DC:3"},
			wantCode:   1,
			wantStderr: "Error: invalid Ethernet address\n",
		},
		{
			name:       "missing address",
			args:       nil,
			wantCode:   2,
			wantStderr: "Usage: macf [-f format] [-u] [ethernet_address]\n\nIf no ethernet_address is provided, macf reads one from stdin.\n\nFormats:\n  cisco   a1b2.c3d4.e5f6\n  colon   a1:b2:c3:d4:e5:f6\n  hyphen  a1-b2-c3-d4-e5-f6\n  none    a1b2c3d4e5f6\n",
		},
		{
			name:       "too many addresses",
			args:       []string{"c6:89:f2:d2:DC:3E", "aa:bb:cc:dd:ee:ff"},
			wantCode:   2,
			wantStderr: "Usage: macf [-f format] [-u] [ethernet_address]\n\nIf no ethernet_address is provided, macf reads one from stdin.\n\nFormats:\n  cisco   a1b2.c3d4.e5f6\n  colon   a1:b2:c3:d4:e5:f6\n  hyphen  a1-b2-c3-d4-e5-f6\n  none    a1b2c3d4e5f6\n",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			stdin := io.Reader(strings.NewReader(tc.stdin))
			gotCode := run(tc.args, stdin, &stdout, &stderr)
			if gotCode != tc.wantCode {
				t.Fatalf("run(%v) code = %d, want %d", tc.args, gotCode, tc.wantCode)
			}
			if stdout.String() != tc.wantStdout {
				t.Fatalf("run(%v) stdout = %q, want %q", tc.args, stdout.String(), tc.wantStdout)
			}
			if stderr.String() != tc.wantStderr {
				t.Fatalf("run(%v) stderr = %q, want %q", tc.args, stderr.String(), tc.wantStderr)
			}
		})
	}
}
