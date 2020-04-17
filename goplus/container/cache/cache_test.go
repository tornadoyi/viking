package cache

import "testing"

func TestLRU(t *testing.T)  {
	c := NewLRU(3)

	c.Add("k1", 1)
	c.Add("k2", 2)
	c.Add("k3", 3)
	c.Add("k4", 4)

	if _, ok := c.Get("k1"); ok { t.Fatal("k1 should be evicted") }
	if _, ok := c.Get("k2"); !ok {t.Fatal("k2 missing")}

	c.Add("k5", 5)
	if _, ok := c.Get("k3"); ok {t.Fatal("k3 should be evicted")}
}
