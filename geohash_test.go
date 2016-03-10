package gwork

import (
	"math"
	"strconv"
	"testing"
)

func TestGeoHash(t *testing.T) {
	smap := NewMap(64800, 64800)
	for i := 1; i < 100; i++ {
		smap.AddCoord("key"+strconv.Itoa(i), i*100, i*100)
	}
	ghash, box := smap.Encode(5000, 5000)
	node0 := &CoordNode{
		x:       5000,
		y:       5000,
		Geohash: ghash,
	}
	cases := smap.QueryNearestSquare(node0.x, node0.y)
	for _, k := range cases {
		node, _ := smap.GetCoordNode(k)
		distance := getDistance(node, node0)
		if distance > int(math.Sqrt(float64(box.Height()*box.Height()+box.Width()*box.Width()))) {
			t.Errorf("key %q not valid", k)
		}
	}
}
