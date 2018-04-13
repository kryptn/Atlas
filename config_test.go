package main

import "testing"

func TestConfig(t *testing.T) {
	c := Config{}
	c.Init()
	c.Load()
	if c.host != "localhost" {
		t.Fail()
	}
}
