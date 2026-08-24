# Codex 移交：文学空间光实现

规格正文：`LITERARY_SPATIAL_LIGHTING.md`（Ready for implementation）。本文件只负责「怎么交给 Codex 做完」，不要改规格里已拍板的数字与产品边界。

工作目录：`E:\blog\astro-blog`

---

## 操作方法

### 1. 打开工程

```powershell
cd E:\blog\astro-blog
```

- **Codex CLI：** 在该目录运行 `codex`（或你本机的 Codex 启动命令），把下面「总提示词」整段贴进去。
- **Codex VS Code / Cursor 侧栏：** 打开 `astro-blog` 文件夹，`@LITERARY_SPATIAL_LIGHTING.md` `@CODEX_HANDOFF.md`，再贴总提示词。

建议 **Agent / 可写文件** 模式，不是 Ask。不要开多个并行 Agent 改同一套 lighting 文件。

### 2. 推荐节奏（一次会话做完全部）

按规格 PR 1 → 8 **顺序提交**，每个 PR 一个 git commit，标题用规格里的英文 title。

```
lighting: persist atmosphere stage and add spatial registers
lighting: split aperture, komorebi, and field without display:none
lighting: rewrite dawn-window as paper-screen aperture
lighting: replace foliage illustration with komorebi
lighting: replace starfield with milky field and dawn halo
motif: share lighting direction in local frame; hide moon in summer
lighting: lock parent opacity at 1 and drop the 50% fade
lighting: About whisper, design language §11, lab copy
```

PR 3 / 4 / 5 在 PR 2 之后可并行，但 **同一会话里仍按 3→4→5 顺序写**，避免抢 `lighting.css`。PR 7 必须在 3–5 之后（盘还在时摘帽会把 blob 调亮）。PR 8 建议在 7 之后。

若上下文不够：第一会话只做 PR 1–2（骨架能走），第二会话 PR 3–5（手艺），第三会话 PR 6–8（motif、摘帽、文学）。下一会话第一句写：

> 继续 `LITERARY_SPATIAL_LIGHTING.md`。先 `git log --oneline -12`，从尚未落地的下一个 PR 接着做，不要重写已合并的 persist / register。

### 3. 每个 commit 后自检

```powershell
npm run check
npm run build
```

开发预览：`npm run dev`，然后走这条路径（验收 13）：

1. `/` Summer → `/music/` → 任意 `/blog/[id]/` → 回 `/`
2. 切 Night 再走一遍
3. `/prototype/acg/lighting/` 点六态，再点「在首页 / 音乐 / 文章检查」
4. 离开预览 URL 后气氛必须回到该页 SSR register
5. 裸 `/?lightingPreset=silver-river` **不得**改正式首页

DevTools 抽检：

- `.environment-lighting` 有 `transition:persist="environment-lighting"`，换页节点不销毁
- 父级 `opacity: 1`，**没有** `.08 / .0515 / .0305` 那六条覆盖（PR 7 之后）
- Home Summer `--lighting-dawn-layer: .22`，热核约 `22vmax`
- Lab 排他根是 `.22` / `.18`，**不是** `1`
- 三根层与朝窗子层 **没有** `display: none`

### 4. 不要做的事

- 不要给访客做三预设选择器
- 不要 `display: none` 切断换房（包括热核 / 尘 / 水 / 梃 2）
- 不要把 mix 再乘到 recipe 上；`lightingRegisterMix` 只是备忘
- 不要再乘 0.5；过吵就缩核或降 local
- 不要 Canvas / WebGL / 粒子 / 视频 / 拟人日月 / 霓虹
- 不要改 `nabunana:environment-v1` 语义（只存 `summer|night`）
- 不要动唱片廊、播放器信息架构、导航 IA
- 不要在开关旁、首页页脚、音乐页、文章正文放诗
- 不要给 product / dev 原型开文学光

---

## 总提示词（整段复制给 Codex）

```text
你在 E:\blog\astro-blog 实现已定稿的文学空间光。规格是唯一真相：先通读 LITERARY_SPATIAL_LIGHTING.md，按文末 PR Plan 1→8 顺序落地，每个 PR 一个 git commit，标题用文档里的英文 title。不要发明第二套产品。

# 产品

站点是一座房子、三层景深，不是三套皮肤，也不是访客换装盘。
- 朝窗 = 室内孔径（纸门，不是太阳）。图纸 full / skeleton / quiet。
- 松影 = 中景木漏日，落在留白和卡片外沿，不画树。
- 星河 = 窗外远场。晨雾+一处日晕；夜是乳白薄雾，几乎无星。
访客控件仍然只有 Summer / Night（nabunana:environment-v1 只存 summer|night）。
页面音域 data-lighting-register：
  / 与 /about → threshold（朝窗 full）
  /blog 与 /projects → grove（松影主导 + 朝窗 skeleton）
  /music → field（星河主导 + 朝窗 skeleton）
  /blog/[id] → quiet（朝窗 quiet + 阅读栅格 mask）
  /prototype/acg/lighting → lab（preset 排他，用 display 行，不是 opacity 1）
  /prototype 与 /prototype/minimal → quiet
  /prototype/product 与 /prototype/dev → 隐藏气氛层
朝窗骨架跨页在场。换页是走进更深的房间，不是换主题站。

已拍板：
1. 触摸天气视差保留 2px（桌面 4px）；prefers-reduced-motion 才归零。
2. About 低语 = 短句 + 作者，无 href。Summer「大きな赤い日」— 夏目漱石；Night「皎皎空中孤月轮」— 张若虚。只出现在 About .place，不是 .about-footer。
3. Home 第一轮就上极淡星河远场：Summer .05 / Night .04。
4. Lab 可以把排他气氛带到首页/音乐/文章，必须 lightingPreview=1；离开 URL 后 after-swap 写回 SSR register。裸 ?lightingPreset= 在正式页无效。
5. WindWaveMotif 轴文字 Night = NIGHT / 2026。

# 合成（唯一公式）

painted = 1 × atmosphereRootOpacity[register, env, intensity] × local × gradient
- .environment-lighting { opacity: 1 } 锁死。PR 7 删除六条 50% 淡化覆盖（Summer .08/.0515/.0305 与 Night 对应值）。
- --lighting-dawn-layer / --lighting-pine-layer / --lighting-river-layer 是已经混合的查表结果。选择器绑在：
  html[data-environment] body[data-lighting-register][data-lighting-intensity]
- 禁止运行时再乘 lightingRegisterMix。禁止 Lab 写成 1/0/0。
- 朝窗子层始终挂在 persist 树上，full→skeleton 用 0.65s opacity，禁止 display:none，禁止 data-dawn-drawing 卸树。

数字地板（Acceptance 12 / 15）：
- Home Summer display 主导层 ∈ .20–.24（表值 .22），热核 22vmax，核有效 ≥ .08，梃有效 ≥ .05
- Night 低于同一格的 Summer
- Music field 上 35% 洗 painted ≥ .05，紫 painted ≤ .08（紫 stop .36 × river .22 = .079）
- 文章 quiet .08 / Night .07；mask 对齐实装栅格 720 + 5rem + 180，居中于 1180 容器（1440 上洞从 ~230px 起，不是视口居中 720）。≤900px 收成 720；≤800px 关洞。

# 实现要点

Persist（PR 1 就必须做）：
- .environment-lighting 使用 transition:persist="environment-lighting"（对标播放器 nabunana-music-player 与开关 environment-controller）
- astro:after-swap 从新 body[data-lighting-register] 和 intensity 抄到 html
- 三根层 transition: opacity 0.65s var(--ease)
- BaseLayout 增加 lightingRegister prop，默认 quiet；inline script 不再写死 dawn-window 为站点身份

拆分：
- src/components/lighting/DawnAperture.astro
- src/components/lighting/PineKomorebi.astro
- src/components/lighting/RiverField.astro
- src/styles/lighting.css
- src/components/lighting/LightingWhisper.astro（PR 8）
EnvironmentLighting.astro 只保留壳、开关、JS。

手艺（坐标在规格里，照抄不要发明第二套）：
- 删 .lighting-source / .lighting-ray / ellipse 树 / 五星 / 扁椭圆 galaxy / stroke-dashoffset 水
- 纸门梃 14% / 39% / 71%，横档 22% / 58%
- 松影 10 斑 + 6 针坐标表；Home ::after 含 .watching-grid article，overflow:visible 在 article 上
- 星河无 stars；Home 远场本轮交付
- Motif：月亮在 motif 本地 % 并 clamp；Summer opacity 0；不要把视口 --lighting-source-x 塞进 SVG cx；Hero 6×4 视差不要叠 --light-shift-*

文学：
- 开关旁零诗句
- 首页保持「また、どこかで。」不另加诗
- 音乐页不加 colophon
- 文章零诗句
- Lab 状态卡：短句衬线 + 作者/作品等宽 + https 链接（杜甫/王维用 Wikisource，不要 CCTV）
- 改写 DESIGN_LANGUAGE.md §11（推翻「正式只有 dawn-window」和 50% 淡化）；§3 Night 底色与代码对齐为 #102724

# 验收句子（必须对照，不允许“差不多”）

1. 看成动画就太快（天气 45–90s，数像素）
2. 页面变橙则 Summer 失败，边缘仍是青绿
3. 正文对比不得掉；字形上无斑
4. 能看出三个径向渐变则手艺失败
5. Night 不能像星空壁纸
6. 不能描出一棵树
7. 不能描出太阳/月亮角色
8. 首页/音乐/文章必须是同一座房子的不同进深
9. 开关旁无诗；About 低语不在 #0a201d 页脚
10. reduced-motion 天气静止但房间仍在
11. 正式站无预设选择器
12. 禁止再乘 0.5
13. Home→Music→文章 0.65s 走，无硬切；persist 必须在
14. 裸 lightingPreset 不得改正式页
15. Music 远场不得被 gradient 打回 2%

每个 PR 后跑 npm run check。视觉 PR 在 /、/music/、/blog/、一篇文章、/prototype/acg/lighting/、About 上检查 Summer 与 Night。不要提交 node_modules 或 dist。
做完后用中文列：改了哪些文件、每个 commit 对应哪条 PR、哪些验收你已用 DevTools 核对、哪些需要我目视。
```

---

## 若只做一个 PR（备用短提示）

把 `N` 换成 1–8：

```text
只实现 LITERARY_SPATIAL_LIGHTING.md 的 PR N，不要做后续 PR。
先读该 PR 的标题、影响文件、依赖、说明，以及 Key Decisions 与 Acceptance。
依赖未落地就先停下来告诉我。完成后 git commit（用文档英文 title），跑 npm run check，用中文汇报。
```

PR 1 额外加一句：视觉仍可走现行 dawn-window CSS，但舞台必须 persist，页面必须写入 register。

PR 7 额外加一句：盘必须已经不在；若 3–5 没删 .lighting-source，先停。
