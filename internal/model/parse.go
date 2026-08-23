package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"bezier-curv/internal/geom"
)

// ParseSpec 从 JSON 字节流解析规格。
//
// 接受两种形态：
//
//	数组： [[x,y],[x,y],[x,y],[x,y]]
//	对象： {"controlPoints": [[x,y],…], "offsetDistance": 0.25}（也可用 "points" 别名）
//
// 控制点既可以是 [x,y] 二元数组，也可以是 {"x":…,"y":…} 对象。
// 解析只做结构转换，合法性由 Validate 检查。
func ParseSpec(data []byte) (Spec, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return Spec{}, errors.New("empty input: no control points")
	}
	if trimmed[0] == '[' {
		var pts []jsonPoint
		if err := json.Unmarshal(trimmed, &pts); err != nil {
			return Spec{}, fmt.Errorf("parse control point array: %w", err)
		}
		return Spec{ControlPoints: toVecs(pts)}, nil
	}
	if trimmed[0] != '{' {
		return Spec{}, errors.New("invalid input: expected JSON array of points or object with controlPoints")
	}
	var js jsonSpec
	if err := json.Unmarshal(trimmed, &js); err != nil {
		return Spec{}, fmt.Errorf("parse spec object: %w", err)
	}
	pts := js.ControlPoints
	if len(pts) == 0 {
		pts = js.Points
	}
	return Spec{
		Name:           js.Name,
		ControlPoints:  toVecs(pts),
		OffsetDistance: js.OffsetDistance,
	}, nil
}

// ParseFile 从文件读取并解析规格。路径为 "-" 时读 stdin（io.Reader）。
func ParseFile(path string, stdin io.Reader) (Spec, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return Spec{}, fmt.Errorf("read input: %w", err)
	}
	return ParseSpec(data)
}

// jsonSpec 是对象形态的 JSON 结构。
type jsonSpec struct {
	Name           string      `json:"name,omitempty"`
	ControlPoints  []jsonPoint `json:"controlPoints,omitempty"`
	Points         []jsonPoint `json:"points,omitempty"`
	OffsetDistance float64     `json:"offsetDistance,omitempty"`
}

// jsonPoint 兼容 [x,y] 与 {"x":…,"y":…} 两种控制点写法。
type jsonPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// UnmarshalJSON 让 jsonPoint 既能吃二元数组也能吃对象。
func (p *jsonPoint) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return errors.New("control point must not be null")
	}
	if trimmed[0] == '[' {
		var pair [2]float64
		if err := json.Unmarshal(trimmed, &pair); err != nil {
			return fmt.Errorf("control point pair: %w", err)
		}
		p.X, p.Y = pair[0], pair[1]
		return nil
	}
	var obj struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := json.Unmarshal(trimmed, &obj); err != nil {
		return fmt.Errorf("control point object: %w", err)
	}
	p.X, p.Y = obj.X, obj.Y
	return nil
}

// toVecs 把解析得到的点列表转成 geom.Vec2 切片。
func toVecs(pts []jsonPoint) []geom.Vec2 {
	out := make([]geom.Vec2, len(pts))
	for i, p := range pts {
		out[i] = geom.Vec2{X: p.X, Y: p.Y}
	}
	return out
}
