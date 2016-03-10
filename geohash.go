package gwork

import (
	"bytes"
	"math"
)

// map 64800
const (
	screenWidth  int = 1920
	screenHeight int = 1080
	maxReturn    int = 40
	BASE32           = "0123456789bcdefghjkmnpqrstuvwxyz"
)

var (
	bits   = []int{16, 8, 4, 2, 1}
	base32 = []byte(BASE32)
)

type CoordNode struct {
	x       int
	y       int
	Geohash string
}

type Box struct {
	MinX, MaxX int
	MinY, MaxY int
}

func (this *Box) Width() int {
	return this.MaxX - this.MinX
}

func (this *Box) Height() int {
	return this.MaxY - this.MinY
}

type MapType struct {
	width          int
	height         int
	precision      int
	allCoordNodes  map[string]*CoordNode
	geohashMapkeys map[string]map[string]bool
}

//创建全局的map
func NewMap(width int, height int) *MapType {
	var grid int
	if screenWidth > screenHeight {
		grid = screenHeight
	} else {
		grid = screenWidth
	}
	precision := int(math.Ceil(float64(height) / float64(grid)))
	precision = 3
	return &MapType{width, height, precision, make(map[string]*CoordNode), make(map[string]map[string]bool)}
}

func (this *MapType) Encode(x, y int) (string, *Box) {
	var geohash bytes.Buffer
	var minX, maxX int = 0, this.width
	var minY, maxY int = 0, this.height
	var precision = this.precision
	var mid int = 0

	bit, ch, length, isEven := 0, 0, 0, true
	for length < precision {
		if isEven {
			if mid = (minY + maxY) / 2; mid < y {
				ch |= bits[bit]
				minY = mid
			} else {
				maxY = mid
			}
		} else {
			if mid = (minX + maxX) / 2; mid < x {
				ch |= bits[bit]
				minX = mid
			} else {
				maxX = mid
			}
		}

		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32[ch])
			length, bit, ch = length+1, 0, 0
		}
	}

	b := &Box{
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
	}
	return geohash.String(), b
}

func (this *MapType) GetNeighbors(x, y int) []string {
	geohashs := make([]string, 9)
	// 本身
	geohash, b := this.Encode(x, y)
	geohashs[0] = geohash

	// 上下左右
	geohashUp, _ := this.Encode((b.MinX+b.MaxX)/2, (b.MinY+b.MaxY)/2+b.Height())
	geohashDown, _ := this.Encode((b.MinX+b.MaxX)/2, (b.MinY+b.MaxY)/2-b.Height())
	geohashLeft, _ := this.Encode((b.MinX+b.MaxX)/2-b.Width(), (b.MinY+b.MaxY)/2)
	geohashRight, _ := this.Encode((b.MinX+b.MaxX)/2+b.Width(), (b.MinY+b.MaxY)/2)

	// 四个角
	geohashLeftUp, _ := this.Encode((b.MinX+b.MaxX)/2-b.Width(), (b.MinY+b.MaxY)/2+b.Height())
	geohashLeftDown, _ := this.Encode((b.MinX+b.MaxX)/2-b.Width(), (b.MinY+b.MaxY)/2-b.Height())
	geohashRightUp, _ := this.Encode((b.MinX+b.MaxX)/2+b.Width(), (b.MinY+b.MaxY)/2+b.Height())
	geohashRightDown, _ := this.Encode((b.MinX+b.MaxX)/2+b.Width(), (b.MinY+b.MaxY)/2-b.Height())

	geohashs[1], geohashs[2], geohashs[3], geohashs[4] = geohashUp, geohashDown, geohashLeft, geohashRight
	geohashs[5], geohashs[6], geohashs[7], geohashs[8] = geohashLeftUp, geohashLeftDown, geohashRightUp, geohashRightDown

	return geohashs
}

func (this *MapType) NewCoordNode(x, y int) *CoordNode {
	ghash, _ := this.Encode(x, y)
	return &CoordNode{
		x:       x,
		y:       y,
		Geohash: ghash,
	}
}

func (this *MapType) GetAllCoordNodes() map[string]*CoordNode {
	return this.allCoordNodes
}

// 增加key的坐标节点
func (this *MapType) AddCoordNode(key string, coordNode *CoordNode) {
	this.allCoordNodes[key] = coordNode

	if this.geohashMapkeys[coordNode.Geohash] == nil {
		this.geohashMapkeys[coordNode.Geohash] = make(map[string]bool)
	}
	this.geohashMapkeys[coordNode.Geohash][key] = true
}

func (this *MapType) AddCoord(key string, x, y int) {
	ghash, _ := this.Encode(x, y)
	coordNode := &CoordNode{
		x:       x,
		y:       y,
		Geohash: ghash,
	}
	this.AddCoordNode(key, coordNode)
}

// 删除key的坐标节点
func (this *MapType) DeleteCoordNode(key string) bool {
	if _, ok := this.allCoordNodes[key]; !ok {
		return false
	}

	ghash := this.allCoordNodes[key].Geohash
	delete(this.geohashMapkeys[ghash], key)
	delete(this.allCoordNodes, key)

	if len(this.geohashMapkeys[ghash]) == 0 {
		this.geohashMapkeys[ghash] = nil
	}

	return true
}

// 更新key的坐标节点
func (this *MapType) UpdateCoordNode(key string, coordNode *CoordNode) bool {
	if !this.DeleteCoordNode(key) {
		return false
	}
	this.AddCoordNode(key, coordNode)
	return true
}

func (this *MapType) UpdateCoord(key string, newx, newy int) bool {
	if !this.DeleteCoordNode(key) {
		return false
	}
	this.AddCoord(key, newx, newy)
	return true
}

// 得到key的坐标节点
func (this *MapType) GetCoordNode(key string) (*CoordNode, bool) {
	coordNode, ok := this.allCoordNodes[key]
	return coordNode, ok
}

// 查找key附近(九宫格内)的节点，返回他们的key
func (this *MapType) QueryNearestSquareFromKey(key string) []string {
	if coordNode, ok := this.GetCoordNode(key); ok {
		return this.QueryNearestSquare(coordNode.x, coordNode.y)
	}
	return []string{}
}

// 查找(x, y)附近(九宫格内)的节点,返回它们的key
func (this *MapType) QueryNearestSquare(x, y int) []string {
	keys := make([]string, 0)
	neighbors := this.GetNeighbors(x, y)
	for _, ghash := range neighbors {
		if this.geohashMapkeys[ghash] != nil {
			for key, _ := range this.geohashMapkeys[ghash] {
				keys = append(keys, key)
			}
		}
	}
	return keys
}

//返回距离的平方
func getDistance(p1 *CoordNode, p2 *CoordNode) int {
	return int(math.Sqrt(float64((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y))))
}
