# nabunana Blog Design Language

> ACG × Music × Editorial — 一个喜欢写代码、听音乐和记录生活的人，为自己搭建的互联网小空间。

## 1. Design Position

这套设计位于“高级音乐／艺术网站”与“可长期维护的个人博客”之间：

- 60% Functional：导航、阅读、搜索、内容归档必须清晰。
- 40% Atmospheric：通过色彩、留白、排版、风与水波营造个人气质。
- ACG 是人格的一部分，不是覆盖全部内容的主题皮肤。
- Yorushika、n-buna、Vocaloid 只作为情绪与音乐文化语境，不复制任何品牌源码、Logo、SVG、插画或独特页面布局。

核心关键词：清澈、克制、安静、透明、轻盈、文学、青春、音乐、夏日、海、风、蓝色时刻、数字花园。

## 2. Signature Motif

原创核心符号为“风与水波”：

- 同一条有机曲线可以被理解为风、海面或音频波形。
- 使用 1px–1.5px 低对比描边。
- 图形只作为 Atmosphere Layer，不抢占正文与信息层级。
- 可辅助使用飞鸟、海平线和唱片圆环；几何月圈只在 Night 的 Motif 本地坐标内以不高于 `.2` 的透明度出现，Summer 隐藏。

对应组件：`src/components/acg/WindWaveMotif.astro`。

## 3. Color System

### Light / Summer

| Token | Value | Usage |
| --- | --- | --- |
| Background | `#F5FBFA` | 页面主背景 |
| Surface | `#FFFFFF` | 轻量内容表面 |
| Secondary Surface | `#EDF7F5` | 分区、抽象封面 |
| Main Text | `#173B3A` | 正文 |
| Strong Text | `#102E2D` | 标题 |
| Muted Text | `#64817F` | 元信息 |
| Primary Accent | `#3DAFA5` | 链接、编号、波形 |
| Deep Accent | `#216E69` | 实心按钮、强调文字 |
| Sky Accent | `#83C9D7` | 天空与水面辅助色 |
| Border | `rgba(39,108,103,.14)` | 细线分层 |

### Dark / Night

| Token | Value | Usage |
| --- | --- | --- |
| Background | `#102724` | 墨绿月夜，不使用纯黑 |
| Surface | `#18332F` | 内容表面 |
| Secondary Surface | `#1E3B36` | 分区 |
| Main Text | `#E7F0ED` | 正文 |
| Strong Text | `#F3F8F6` | 标题 |
| Muted Text | `#9EB3AE` | 元信息 |
| Accent | `#78C5BB` | 月光青色 |
| Soft Accent | `#91D6CD` | 低频强调 |
| Border | `rgba(174,222,215,.15)` | 深色细线 |

粉、紫、黄只能用于很小的状态或内容点缀，不得成为大面积背景。

## 4. Typography

- 大标题使用衬线字体，强调文学性、呼吸感和中日文字形。
- 正文使用系统无衬线字体，确保中文长文稳定可读。
- 编号、日期、标签使用等宽字体。
- 日文只承担氛围与辅助语义，中文或英文必须保证可理解性。
- 标题允许紧凑字距；正文禁止过度压缩。

推荐层级：

- Hero：`clamp(3.8rem, 7.6vw, 8.2rem)`
- Section title：`clamp(2rem, 4.2vw, 4.6rem)`
- Article title：`clamp(2.8rem, 7vw, 5.8rem)`
- Reading body：约 `1.05rem / 1.85`，最大宽度 `720px`

## 5. Layout

- 页面容器上限约 `1180px`。
- 主要依靠留白、细线和明暗对比建立层级。
- 首页仍按 Hero → Currently → Writing → Music → Notes → Projects → Watching → Fragments → About 排列，不重排信息架构。Hero、Currently、Music、Watching 是四幕 scroll-driven 电影章；Writing、Notes、Projects、Fragments、About 是五段普通编辑流。电影章之间由文档流留出呼吸，页面不是 scroll-snap 式全屏演示。
- 电影章使用 `track + sticky pin + stage`：桌面 pin 固定在 `top: 0`、高 `100svh`，底部至少留出 `80px` 给播放器与环境控件；锚点始终落在 track 顶。站点 Header 保持普通文档流。
- 文章列表使用 Editorial List，不使用重复圆角卡片。
- 正文阅读优先，ACG 视觉不得侵入文章主体。
- 图片允许 16:9、4:3、3:4、2:3 和非对称裁切；真实专辑封面保持 1:1。

## 6. Component Rules

- 默认圆角：`6px`；轻量容器可用 `10px`；少量媒体内容上限 `14px`。
- 默认无阴影；Hover 阴影最多 `0 12px 40px rgba(25,91,86,.08)`。
- Glassmorphism 只用于导航、搜索、浮动 Summer/Night 控件和全站迷你播放器。
- 按钮数量要少，优先文本链接与箭头反馈。
- Hover 可使用 `translateX(4px)`、下划线展开或 2px 轻微上浮。
- 禁止 Glow、霓虹、厚重 Dashboard 和游戏启动器视觉。

## 7. Music Identity

- 播放控件只对应真实、可访问的音源，不显示虚构进度或播放时长。
- 全站播放器固定在左下角，采用深夜海绿、低对比边框和轻量模糊；默认收为 60px 封面，悬停、聚焦或点按后展开，离开焦点约 360ms 后以约 200ms 动画快速收回。用户可手动固定展开状态，固定选择随浏览器保存；手动收起同时解除固定。播放器不能遮挡正文、Mood 控件或移动端导航。
- 播放器默认不自动播放，保留用户选择的歌曲、位置和音量；专辑封面可从首页直接开始播放。
- 展开播放器使用与播放器主体同宽的独立歌词层显示日文原文与同句中文译文；根据整首 LRC 的日中交替结构配对同句翻译，不能把译文与下一句原文错配。纯音乐明确显示 `Instrumental / 无歌词`。歌曲选择面板位于歌词上方，按专辑折叠分组，任何视口下都不得与歌词重叠。
- 站内导航使用 Astro `ClientRouter`，播放器根节点必须保持同名 `transition:persist`，禁止在页面切换时重建正在播放的 `<audio>`。
- 当前音乐库为七张用户确认的专辑，共 94 首；首页暗房章节使用滚动进度连续驱动、不可循环的 3D 封面景深，第一张到最后一张严格对应页面时间轴。水平拖拽和方向键只改变同一页面进度，不设自动巡航、不抢占竖直滚轮；完整专辑资料与封面画廊集中在音乐页。
- Favorite Artists 只展示用户已经确认的 Yorushika 与 n-buna；其他条目使用明确占位。
- Vocaloid 采用 Japanese Music Archive 语义，围绕制作人、作品、听感和个人笔记组织。
- 音乐页允许更深的海绿色背景与唱片抽象图形，但仍需保持博客导航和可读性。

## 8. Motion

- 近景风、水波 Motif 周期为 8s–20s；全站环境天气周期为 45s–90s，层与层不得齐步。
- 天气视差桌面控制在 4px，触摸保留 2px；Motif 自身的 6px × 4px 视差不叠加天气变量。
- 页面反馈使用 180ms–320ms 的淡入、轻微位移或 Blur-lite。
- 首页滚动叙事只有一个 rAF 导演，按 `p = clamp(-track.top / (track.height - viewport), 0, 1)` 写入章节 `--p`；滚动本身零 duration，正向与反向必须沿同一时间轴连续还原。禁止一次性 fade-in、`scroll-snap`、GSAP/Lenis 与滚轮劫持。
- 滚动驱动的文案位移不超过 `32px`，Motif 不超过 `24px`，比例限制在 `0.94–1.04`，Music `rotateY` 不超过 `12deg`。文档流整段只允许 `opacity: .35 → 1` 与 `translateY: 16px → 0`，不得让每条列表分别飞入。
- 禁止 Bounce、高频 Parallax 和快速位移。
- `≤900px` 取消 sticky 加长，章节回到正常文档流；Music 仅保留水平拖拽与键盘切换，不增加巡航。`prefers-reduced-motion` 下 track 塌回内容高度、章节进度冻结，所有内容不依赖动画也能完整阅读。
- 触摸移动端保留 2px 天气视差；只有 `prefers-reduced-motion` 才归零。

## 9. Responsive Rules

- `1440px`：完整双栏 Hero、完整导航与多栏内容。
- `768px`：Hero 转单栏，内容网格降级，保留足够留白。
- `430px`：长标题收缩、Watching 单列、Fragments 纵向排列。
- `375px`：使用与 430px 相同的移动规则，按钮和控件保持至少 44px 触控区域。
- 固定 Summer/Night 控件与迷你播放器不得互相遮挡，也不得遮挡原型工具条或正文。

## 10. Content Policy

- 正式站只读取 `src/content/blog`、`src/content/projects` 和 `src/content/notes`。
- 不导入旧 Hexo 文章；旧站内容可以保留为本地历史备份，但不进入 Astro 构建和导航。
- 不虚构用户喜欢的歌手、专辑、演出、照片或听歌历史。
- Photography、Live 等没有真实内容时，应显示克制的待填充状态；Playlist 使用 `src/data/music.ts` 中的真实专辑和曲目数据。
- ACG Watching 卡片允许使用作品官网公开的宣传视觉图，但必须本地优化、提供准确替代文本并链接回官方来源，不使用第三方热链或来源不明的同人图。
- 文章、项目和状态文字优先使用真实数据。

## 11. Accessibility & Performance

- 所有交互必须支持键盘焦点和语义化标签。
- 搜索支持 `Cmd/Ctrl + K`。
- Summer / Night 是全站唯一的亮暗主题状态：首次访问跟随系统，手动选择后写入 `nabunana:environment-v1`。`data-theme` 只作为旧样式的派生兼容值，不得存在第二套独立主题控制器。
- 使用 Astro 静态输出、原生 CSS 与少量原生 JavaScript。
- 装饰 SVG 应内联、低复杂度且标记为非正文内容。
- 保留 canonical、Open Graph、RSS、sitemap、robots 和文章 JSON-LD。

### Literary Environment Lighting

- 正式站采用已验收的「星河」作为唯一文学空间光：Summer 是淡紫晨雾与一处暖晕，Night 是乳白薄雾与墨绿月夜；访客仍只切换 Summer / Night，不提供预设换装盘。
- 环境状态写在 `<html data-environment="summer|night" data-lighting-preset="silver-river">`；页面保留 `<body data-lighting-register="threshold|grove|field|quiet|lab" data-lighting-intensity="display|standard|reading">`，register 继续表达房间位置，正式站的可见振幅由 intensity 决定。
- 正式页 display / standard / reading 的星河根强度分别为 Summer `.22 / .14 / .08`、Night `.18 / .12 / .07`。三根光影层仍挂在 persist 舞台上，以 0.65s opacity 行走，禁止 `display:none`；朝窗与松影保留给 Lab 对照，不参与正式页合成。
- 唯一合成公式为 `painted = 1 × --lighting-*-layer × local × gradient`。`.environment-lighting` 父级 opacity 锁死为 1；层 token 已经混合 register、environment 与 intensity，禁止再乘 mix 或 0.5。
- Home 与 Music 的星河 display 层为 Summer `.22` / Night `.18`，紫 stop `.36`；文章 reading 为 Summer `.08` / Night `.07`。Night 同格强度始终低于 Summer。
- 朝窗使用 14% / 39% / 71% 竖梃和 22% / 58% 横档，Summer 保留小面积暖核，Night 使用横向银洗与底缘有机水纹；不直接绘制太阳或环境光月亮。
- 松影只画 10 个疏斑与 6 条夜间银针，不画树；星河只画晨雾、单一日晕与乳白薄雾，不画星空壁纸。
- 正式日夜控件只显示 Summer / Night，不展示诗句、出处或额外说明卡片。
- 文学短句只在 Lighting Lab 以全文与来源出现，并在 About `.place` 留一句“短句 + 作者”的无链接低语；首页、音乐页、文章、开关与深色页脚均为零。
- 同一光源方向必须影响背景、线性表面高光、边缘和 Motif 的弱投影；禁止黄色/蓝色全屏滤镜、重 Bloom、视频背景、大贴图和满卡径向高光。
- 文章使用对齐实装 `720px + 5rem + 180px` 栅格的阅读 mask，1440px 上洞从约 230px 开始；≤900px 收为 720px，≤800px 关闭。任何光影不得降低正文对比度。
- Night 正式采用“墨绿月夜”：全站底色为 `#102724`，表面使用 `#18332f` / `#1e3b36`，文章页使用略浅的 `#122a27` 阅读底色，不使用中性纯黑。
- 环境天气只使用原生 CSS、低复杂度内联 SVG 和少量 Pointer Events；周期 45–90 秒、位移数像素，页面隐藏时暂停。`prefers-reduced-motion` 下天气静止，但静态房间仍在。
- `.environment-lighting` 使用 `transition:persist="environment-lighting"`；Home → Music → 文章通过三根层的 0.65s opacity 交叉，不依赖整页溶解。
- 开发比较器位于 `/prototype/acg/lighting/`，以 display 行排他比较三套手艺 × 两种环境；链接只有带 `lightingPreview=1` 才能临时覆盖正式站的星河状态，裸 `lightingPreset` 无效且选择不写入正式偏好。

## 12. Avoid

- 粉色萌系二次元模板
- 大面积紫色渐变或 RGB Neon
- 动漫人物占据主要 Hero
- Spotify、VS Code、Discord 或游戏启动器式复刻
- 每个容器都做卡片、玻璃和阴影
- 版权插画作为默认素材
- Yorushika 官网或 n-buna artwork 的直接复制
- 为填满页面而虚构内容
- 把三层文学空间光做成正式站预设选择器，或再次对气氛父级乘 0.5

## 13. Review Question

加入任何视觉或内容元素前，先回答：

> 它是否让这个网站更像这个人，同时没有降低阅读与导航质量？

如果答案是否定的，就删除。
