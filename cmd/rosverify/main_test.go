package main

import (
	"reflect"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	cases := []struct {
		in   string
		path string
		verb string
		ok   bool
	}{
		{"/ip/hotspot/active/print", "/ip/hotspot/active", "print", true},
		{"/interface/monitor-traffic", "/interface", "monitor-traffic", true},
		{"/tool/ping", "/tool", "ping", true},
		{"/ping", "", "ping", true},
		{"ping", "", "ping", true},
		{"/", "", "", false},
		{"", "", "", false},
		{"/ip/", "", "", false},
	}
	for _, c := range cases {
		path, verb, ok := splitCommand(c.in)
		if path != c.path || verb != c.verb || ok != c.ok {
			t.Errorf("splitCommand(%q) = (%q, %q, %v), mau (%q, %q, %v)",
				c.in, path, verb, ok, c.path, c.verb, c.ok)
		}
	}
}

func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"ip/hotspot/active":   "ip,hotspot,active",
		"/ip/hotspot/active":  "ip,hotspot,active",
		"ip, hotspot,,active": "ip,hotspot,active",
		"":                    "",
	}
	for in, want := range cases {
		if got := normalizePath(in); got != want {
			t.Errorf("normalizePath(%q) = %q, mau %q", in, got, want)
		}
	}
}

func TestSummarize(t *testing.T) {
	checks := []Check{
		{Resource: "a", Status: statusOK},
		{Resource: "a", Status: statusWarn},
		{Resource: "b", Status: statusDrift},
		{Resource: "b", Status: statusUnknown},
		{Resource: "b", Status: statusError},
	}
	got := summarize(checks)
	want := Summary{Resources: 2, Checks: 5, OK: 1, Warn: 1, Drift: 1, Unknown: 1, Error: 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("summarize = %+v, mau %+v", got, want)
	}
}
