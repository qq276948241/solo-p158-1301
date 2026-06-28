# 双人贪吃蛇 架构说明

> 用大白话讲清楚这个项目是怎么搭的，没看过代码也能看懂。

---

## 一句话总览

这是一个用 **Go + Bubble Tea** 写的终端双人贪吃蛇游戏。整个程序跑在 Bubble Tea 的 **事件循环**里：每隔 120ms 推进一帧（tick），每条蛇走一步，然后检查吃食物、撞墙、撞自己、撞对方，最后把最新画面画到终端上。玩家按键盘就是在给这个循环发消息。

---

## 整体运行模式：Bubble Tea 的 MVU

Bubble Tea 用的是类似 Elm 的 **Model-View-Update** 模式：

```
┌──────────────┐     消息      ┌──────────────┐     新状态     ┌──────────────┐
│   外部输入    │ ────────────> │   Update()   │ ────────────> │    Model     │
│ (键盘/定时器) │               │              │               │  (游戏状态)   │
└──────────────┘               └──────────────┘               └──────┬───────┘
                                                                     │
                                                                     │ 用Model画
                                                                     ▼
                                                               ┌──────────────┐
                                                               │    View()    │
                                                               │  (渲染字符串) │
                                                               └──────────────┘
```

在咱们项目里：

- **Model** = `UI` 结构体（里面包着 `Game`、`HighScore`），定义在 [ui.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/ui.go#L26-L29)
- **Update** = [UI.Update()](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/ui.go#L245-L291)，处理按键和 tick 消息
- **View** = [UI.View()](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/ui.go#L198-L224)，把游戏状态拼成带颜色的字符串

游戏的主驱动（真正的游戏逻辑）在 **`game.go` 的 `Game.Tick()`** 里。每 120ms Bubble Tea 发一个 `tickMsg`，`Update` 收到后就调一次 `Game.Tick()`，推进一步。

---

## 每个文件是干什么的

### 🏁 [main.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/main.go) — 程序入口

干三件事：

1. **读高分存档**：启动时从同目录的 `score.txt` 读历史最高分，文件不存在就当 0 分。
2. **启动 Bubble Tea**：`NewGame()` 建游戏、`NewUI()` 包成 Model，然后 `tea.NewProgram(ui).Run()` 进入事件循环。
3. **存高分存档**：玩家退出游戏（按 Q 或 Ctrl+C）后，比较 P1/P2 当前分数和历史最高分，把更高的值写回 `score.txt`。

> `score.txt` 就放在可执行文件所在的目录里，内容就是一个纯数字。

---

### 🎯 [game.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/game.go) — 游戏总调度

这是整个项目的"大脑"，把所有模块串起来。

核心数据结构是 `Game` 结构体，里面装着：
- 一张棋盘 `Board`
- 两条蛇 `Snake1`、`Snake2`
- 食物管理器 `FoodMgr`
- 当前状态 `State`（进行中 / 暂停 / 游戏结束）
- 胜者 `Winner`

最关键的方法是 **`Game.Tick()`**，每 120ms 被调用一次，内部按这个顺序执行：

```
1. Snake1.Move()    ← 玩家1走一步
2. Snake2.Move()    ← 玩家2走一步
3. checkFoodCollision(Snake1)   ← 检查玩家1头撞到食物没
4. checkFoodCollision(Snake2)   ← 检查玩家2头撞到食物没
5. 检查各种死亡（撞墙 / 撞自己 / 撞对方 / 吃毒药死）
6. 有人死了就切到 GameOver 状态，判定胜者
7. 累计 25 个 tick（≈3秒）就刷新所有食物
```

还提供 `Reset()`（按 R 重开）和 `TogglePause()`（按空格暂停）两个辅助方法。

---

### 🐍 [snake.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/snake.go) — 蛇的一切

蛇的所有行为都在这里。一条蛇就是：

```go
type Snake struct {
    Body      []Point      // 身体坐标切片，第0个是头
    Dir / NextDir          // 当前方向 / 下个 tick 的方向
    Color                  // 颜色
    Alive                  // 还活着吗
    Score                  // 分数
    GrowQueue              // 待生长的节数（吃了食物先记账，移动时才真正变长）
}
```

主要方法：

| 方法 | 作用 |
|---|---|
| `Move()` | 按方向算出新头，加到 Body 最前面；如果 GrowQueue 没东西就把尾巴切掉一节 |
| `SetDirection()` | 改方向，会拦掉 180° 反向（比如向右时不能直接向左） |
| `Grow()` | `GrowQueue++`，下次 Move 就多长一节 |
| `Shrink(n)` | 从尾巴砍 n 节。砍完会没命就返回 `false`，调用方再决定置 `Alive=false` |
| `CollidesWithSelf()` | 头是否撞到自己身体 |
| `CollidesWith(other)` | 自己的头是否撞到对方蛇的身体 |

两条蛇的初始位置：P1（绿）在左上向右走，P2（蓝）在右下向左走，各 3 节。

---

### 🍎 [food.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/food.go) — 食物系统

用 **配置表 + 效果函数** 的方式管理三种食物：

```
foodConfigs 配置表
├── FoodNormal  🔴 红色  75%  +10分，长1节
├── FoodGolden  🟡 金色  10%  +50分，长2节
└── FoodPoison  🟣 紫色  15%  -20分，缩2节（≤3节直接死）
```

每种食物的效果是一个独立的函数 `FoodEffect func(*Snake)`，配置表里绑定了类型、颜色、权重、效果函数。以后加新食物 **只需要加一个效果函数 + 一条配置**，核心逻辑不用改。

`FoodManager` 负责：

| 方法 | 作用 |
|---|---|
| `Spawn()` | 随机找一个不与蛇/现有食物重叠的格子，按权重随机抽类型生成食物 |
| `SpawnAll()` | 把场上食物补到最多 5 个，并且**兜底保证至少有 1 个普通食物**（不然全是毒药游戏就没法玩了） |
| `Refresh()` | 每 3 秒调用一次，清空所有食物重新生成 |
| `GetFoodAt(p)` | 按坐标取食物引用，蛇吃到时调 `food.ApplyEffect(snake)` 触发效果 |

---

### 🔲 [board.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/board.go) — 棋盘

棋盘是固定 20×20 的。主要就两件事：

1. **边界检测** — `InBounds(p)` 判断一个坐标是否还在棋盘里，撞墙就死。
2. **局面渲染** — `Render(s1, s2, foodMgr)` 把当前两条蛇 + 食物的位置翻译成一个二维 `Cell` 网格，每个格子标清楚是"空 / P1头 / P1身 / P2头 / P2身 / 食物"，UI 层拿到这个网格就只管画颜色，不用关心游戏逻辑。

---

### 🎨 [ui.go](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/ui.go) — 终端渲染

纯展示层，用 `lipgloss` 给字符串上色。整体布局：

```
┌──────────────────────────────┐  ┌─────────────────────────┐
│                              │  │    SNAKE BATTLE         │
│                              │  │                         │
│        20×20 棋盘区          │  │  P1 (WASD)              │
│    (边框 + 蛇 + 食物)        │  │    Score: 20            │
│                              │  │    Length: 5            │
│                              │  │                         │
│                              │  │  P2 (Arrows)            │
│                              │  │    Score: 30            │
│                              │  │    Length: 6            │
│                              │  │                         │
│                              │  │  FOOD                   │
│                              │  │    Count: 3 / 5         │
│                              │  │    Refresh: 3s          │
│                              │  │                         │
│                              │  │  HIGH SCORE             │
│                              │  │    Best: 150            │
└──────────────────────────────┘  └─────────────────────────┘

  🐍 P1: WASD  |  P2: ↑↓←→  |  SPACE: Pause  |  R: Restart  |  Q: Quit
```

- **棋盘区**：从 `Board.Render()` 拿网格，每个格子根据类型选字符（头=●、身=■、食物=◆、空=·）和颜色。
- **右侧面板**：实时显示双方分数、蛇身长度、食物数量、历史最高分、按键说明。
- **浮层**：暂停和游戏结束时在棋盘区中央盖一个提示框。

`UI.Update()` 负责把键盘消息翻译成蛇的方向变化或游戏状态变化（暂停/重开/退出）。

---

## 数据流向（完整链路）

以「玩家 1 按 W，然后吃到金色食物」为例，走一遍完整链路：

```
 玩家按 W 键
    │
    ▼
 tea.KeyMsg 消息进入 UI.Update()
    │
    ├─> 识别到 'w' 键
    └─> Snake1.SetDirection(DirUp)  ← 记录下个 tick 要往上走
    │
    ▼
 120ms 到了，tea.Tick 发 tickMsg
    │
    ▼
 UI.Update() 收到 tickMsg，调 Game.Tick()
    │
    ├─> Snake1.Move()           蛇往上挪一格
    ├─> Snake2.Move()           另一条蛇也挪一格
    │
    ├─> checkFoodCollision(Snake1)
    │      │
    │      ├─> FoodMgr.GetFoodAt(蛇头位置)  ← 发现是金色食物
    │      ├─> FoodMgr.RemoveAt(...)        把这个食物从场上删掉
    │      ├─> food.ApplyEffect(Snake1)
    │      │       └─> goldenEffect()
    │      │             ├─> Snake1.Grow() × 2   GrowQueue += 2
    │      │             └─> Snake1.Score += 50
    │      └─> FoodMgr.Spawn()                补一个新食物
    │
    ├─> 检查撞墙 / 撞自己 / 撞对方  ← 没撞，活着
    │
    └─> 累计 tick 计数（够25次就刷新所有食物）
    │
    ▼
 UI.View() 根据最新 Game 状态重新渲染整个画面
    │
    ├─> Board.Render() 生成 20×20 网格
    ├─> lipgloss 给每个格子上色
    ├─> 拼接右侧分数面板
    └─> Bubble Tea 把字符串刷到终端屏幕上
```

---

## score.txt 的读写时机

| 时机 | 操作 | 代码位置 |
|---|---|---|
| 游戏启动 | 读取历史最高分 | [main.go → loadHighScore()](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/main.go#L16-L37) |
| 游戏退出 | 取 P1/P2 中较高分，若 > 历史最高分则写入 | [main.go → main() 尾部](file:///d:/code/ai-prompt/solo-chrome-dev-F12/repos/repo158/project158/main.go#L71-L84) |

文件位置：**可执行文件所在的目录**下的 `score.txt`（不是源码目录，这样编译好的 exe 放到哪都能正确读写）。内容就是一行纯数字，用记事本打开就能改。
