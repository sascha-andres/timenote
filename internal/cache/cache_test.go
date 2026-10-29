package cache

import "testing"

func TestSetGet(t *testing.T) {
	c, err := NewCache(10, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !c.NeedUpdate("items") {
		t.Error("NeedUpdate of missing key = false, want true")
	}
	if got, err := Get[[]string](c, "items"); err != nil || got != nil {
		t.Errorf("Get of missing key = %v, %v; want nil, nil", got, err)
	}
	if err := Set(c, "items", []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	got, err := Get[[]string](c, "items")
	if err != nil || len(got) != 2 || got[1] != "b" {
		t.Errorf("Get = %v, %v; want [a b]", got, err)
	}
	if c.NeedUpdate("items") {
		t.Error("NeedUpdate after Set = true, want false")
	}
}

func TestNeedUpdateExpired(t *testing.T) {
	c, err := NewCache(0, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := Set(c, "items", 1); err != nil {
		t.Fatal(err)
	}
	if !c.NeedUpdate("items") {
		t.Error("NeedUpdate with max age 0 = false, want true")
	}
}
