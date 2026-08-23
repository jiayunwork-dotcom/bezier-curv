package cli

import (
	"fmt"
	"io"
)

// writeUsage 把用法说明写到 w。
func writeUsage(w io.Writer) {
	fmt.Fprint(w, `bezier-curv — 三次 Bézier 弧长、曲率与法向等距偏移核算工具

用法：
  bezier-curv <command> [flags] <curve.json>

命令：
  sample        打印弧长、若干 t 的曲率 κ(t) 与法向偏移折线
  arclength     弧长细分收敛明细（自适应 Simpson + Gauss 交叉 + 弦长对照）
  curvature     指定参数求曲率（-t 可重复）
  offset        输出法向等距偏移折线（-d 距离，-n 段数）
  validate      只校验控制点是否合法（不合法则 stderr + 退出码 1）
  check         交叉不变性核算（平移/缩放/d=0/共线）
  version       输出版本
  help          显示本帮助

输入：四个控制点的 JSON，文件路径或 "-"（stdin）。
  数组形态：  [[x,y],[x,y],[x,y],[x,y]]
  对象形态：  {"controlPoints":[[x,y],…],"offsetDistance":0.25}

算例：
  bezier-curv sample example/cubic-bow.json
  bezier-curv curvature example/cubic-bow.json -t 0.5
  bezier-curv offset example/quarter-circle.json -d 0.1 -n 16
`)
}

// commandUsage 打印某个子命令的参数用法到 stderr。
func commandUsage(w io.Writer, name, text string) {
	fmt.Fprintf(w, "bezier-curv %s: %s\n", name, text)
	fmt.Fprintf(w, "用法：bezier-curv %s [flags] <curve.json>\n", name)
}
