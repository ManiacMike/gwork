package gwork

// import (
//   "math"
// )

type Point struct{
  x uint
  y uint
}

type MapType struct{
  width uint
  height uint
  points map[string]Point
}

var (
  TheMap *MapType
)

const(
  percise uint = 5
  maxReturn uint = 40
)


//创建全局的map
func NewMap(width uint, height uint) *MapType{
  TheMap = &MapType{width,height,make(map[string]Point)}
  return TheMap
}

//返回离p点最近
func (m *MapType)getNearest(myPoint Point) map[string]Point{

}

//返回距离的平方
func getDistanceSquare(p1 Point,p2 Point) uint{
  return (p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)
}
