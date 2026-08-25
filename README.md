# bezier-curv

三次 Bézier 弧长、曲率与法向等距偏移核算命令行工具。

输入四个控制点（JSON），`bezier-curv` 用 Bernstein 基与数值积分求出：弧长 `∫₀¹|r'(t)|dt`（自适应 Simpson，按「加密一档后变化小于容差」收敛）、曲率 `κ(t)=|r'×r''|/|r'|³`、以及法向等距偏移折线 `r(t)+n̂·d`。能力边界：只做参数曲线核算，不是矢量画板；只处理三次 Bézier，不含凸轮升程/啮合线等其它曲线议题。

## 安装与运行

需要 Go 1.21+。无第三方依赖。

```bash
go build ./...
go run . sample example/cubic-bow.json
```

`sample` 打印弧长、11 个参数点的 κ 与偏移折线：

```text
arc length (adaptive Simpson): 3.091544
  chord |P3-P0|:          2.000000     # 弓形：弧长 > 弦长
  control polygon perm.:  4.041088     # 注意：控制多边形周长不是弧长
```

## 子命令

```
bezier-curv <command> [flags] <curve.json>
  sample      弧长 + 若干 t 的 κ(t) + 偏移折线
  arclength   弧长收敛明细（自适应 Simpson / Gauss 交叉 / 弦长对照）
  curvature   指定参数求曲率（-t 可重复；尖点处输出 err 并退出码 1）
  offset      法向等距偏移折线（-d 距离，-n 段数）
  validate    只校验控制点；不合法 → stderr + 退出码 1
  check       交叉不变性核算（平移/缩放/d=0/共线/弧长≥弦长）
  version / help
```

输入文件路径或 `-`（stdin）。控制点支持数组 `[[x,y],…]` 或对象 `{"controlPoints":[[x,y],…],"offsetDistance":0.25}` 两种形态，单个点可写作 `[x,y]` 或 `{"x":…,"y":…}`。

### 算例

- `example/cubic-bow.json` — 弓形，弧长应大于弦长。
- `example/quarter-circle.json` — 半径 1 的四分之一圆弧（弧长约 1.5710）。
- `example/straight-line.json` — 共线控制点，κ=0，偏移为平行线。

```bash
go run . arclength example/quarter-circle.json
go run . curvature example/cubic-bow.json -t 0.5 -t 0.9
go run . offset example/quarter-circle.json -d 0.1 -n 16
go run . check example/cubic-bow.json
```

## 关键约定

- **弧长**：自适应 Simpson，相对容差 1e-9；`arclength` 命令同时给出 Gauss-Legendre 16 点交叉值与粗弦长和对照。
- **曲率**：`κ=|r'×r''|/|r'|³`，叉积为二维有向标量，分母是速率立方。
- **偏移**：偏移点 = `r(t)+n̂·d`，`n̂` 为左法向；直线段上偏移量恰为 `|d|`，`d=0` 时与原曲线重合。
- **错误可见**：控制点不是 4 个、坐标非有限、曲线退化成一个点 → 校验错误；重合导致的尖点（`r'=0`）处求曲率/偏移 → 明确报错，不返回静默数值。所有非法输入写 stderr 并带非零退出码。

## 交叉不变性

`check` 子命令内置核算：平移所有控制点弧长与曲率不变；整条曲线缩放 `k` 时弧长 ×`|k|`、曲率 ×`1/|k|`；共线控制点 `κ=0` 且偏移为平行线。

## 测试

```bash
go test ./...
```

## 项目结构

```
main.go                  # 子命令接线
internal/geom/           # 二维向量、折线、数值比较
internal/curve/          # Bernstein 基、Bézier 求值/导数/曲率/切法向/尖点检测
internal/arclen/         # 自适应 Simpson、Gauss-Legendre、弦长近似、收敛
internal/offset/         # 法向等距偏移与平行校验
internal/model/          # JSON 解析、校验、尖点诊断、交叉核算
internal/cli/            # 子命令与输出
example/                 # 离线算例
```
