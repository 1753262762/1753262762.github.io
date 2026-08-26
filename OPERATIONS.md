# nabunana 博客编译与维护手册

这份文档用于维护 `E:\blog\astro-blog` 中的 Astro 静态博客。博客文章、页面前端、音乐播放器和静态资源会由同一次构建生成，不需要再运行旧 Hexo 项目。

## 1. 环境与目录

推荐环境：

- Windows 10/11
- Node.js 22.19 或 24
- npm 10 或更高版本
- PowerShell

当前验证环境为 Node.js `v24.19.0`、npm `11.17.0`。

主要目录：

| 路径 | 用途 |
| --- | --- |
| `src/content/blog` | 正式博客文章，Markdown 格式 |
| `src/content/projects` | 项目条目 |
| `src/content/notes` | 首页短笔记 |
| `src/components` | 页面与播放器组件 |
| `src/data/music.ts` | 专辑、歌曲、歌词和播放索引数据 |
| `src/lib/environment.ts` | Summer / Night 环境状态与“星河”正式光影方案 |
| `public/media` | 音频、歌词、封面和图片等静态资源 |
| `dist` | 正式构建产物，可交给 Nginx 提供访问 |

`dist`、`.astro` 和 `node_modules` 都是生成目录，不要在里面直接修改文件。

## 2. 第一次安装

打开 PowerShell：

```powershell
Set-Location -LiteralPath 'E:\blog\astro-blog'
node --version
npm --version
npm ci
```

`npm ci` 会严格按照 `package-lock.json` 安装依赖。只有主动升级依赖时才使用 `npm install`。

## 3. 本地编辑与实时预览

进入项目并启动开发服务器：

```powershell
Set-Location -LiteralPath 'E:\blog\astro-blog'
$env:ASTRO_TELEMETRY_DISABLED = '1'
npm run dev -- --host 127.0.0.1
```

浏览器打开：

- 首页：`http://127.0.0.1:4321/`
- 文章列表：`http://127.0.0.1:4321/blog/`
- 音乐页：`http://127.0.0.1:4321/music/`
- 光照实验：`http://127.0.0.1:4321/prototype/acg/lighting/`

开发服务器运行时，修改 `src` 或 `public` 中的文件通常会自动刷新。按 `Ctrl+C` 停止服务器。

如果 4321 端口被占用：

```powershell
npm run dev -- --host 127.0.0.1 --port 4322
```

## 4. 新增或修改博客文章

在 `src/content/blog` 中创建 UTF-8 编码的 Markdown 文件。文件名会成为文章地址的一部分，建议只用小写英文、数字和连字符，例如：

```text
src/content/blog/my-new-post.md
```

文章模板：

```markdown
---
title: "文章标题"
description: "用于列表、搜索和分享卡片的简短摘要。"
published: 2026-08-20T18:00:00+08:00
updated: 2026-08-20T18:00:00+08:00
draft: true
category: Notes
tags: [Astro, Web]
featured: false
author: nabunana
readingTime: "6 min"
---

从这里开始写正文。
```

注意事项：

- `draft: true` 的文章不会进入正式首页、文章列表、详情页和 RSS。
- 准备公开时改为 `draft: false`，或者删除 `draft` 字段。
- `published` 和 `updated` 使用带时区的 ISO 时间。
- `tags` 必须是数组。
- 不要把旧 Hexo 文章目录重新复制进来；正式站只读取当前 Astro 内容集合。

保存后先在开发服务器中检查文章地址：

```text
http://127.0.0.1:4321/blog/my-new-post/
```

## 5. 音乐资源维护

音乐库的来源文件是：

```text
src/data/music.ts
```

每首歌曲需要对应的音频和歌词：

```text
public/media/music/<album-slug>/01.mp3
public/media/music/<album-slug>/01.lrc
```

专辑封面放在：

```text
public/media/music/covers/<album-slug>.webp
```

维护规则：

- 音频继续使用 160kbps MP3。
- 歌词使用 UTF-8 LRC。
- 文件编号使用两位数字，并与 `music.ts` 中的曲目顺序一致。
- 新专辑追加到现有专辑之后，避免破坏浏览器保存的歌曲索引。
- 每次改动后检查专辑数、歌曲数、MP3 数和 LRC 数是否一致。

## 6. 类型检查与正式构建

只检查 Astro、内容结构和 TypeScript：

```powershell
Set-Location -LiteralPath 'E:\blog\astro-blog'
$env:ASTRO_TELEMETRY_DISABLED = '1'
npm run check
```

生成正式网站：

```powershell
Set-Location -LiteralPath 'E:\blog\astro-blog'
$env:ASTRO_TELEMETRY_DISABLED = '1'
npm run build
```

`npm run build` 会先执行 `astro check`，检查通过后再生成静态网站。成功时应看到：

```text
0 errors
0 warnings
dist/
```

最终需要发布的是整个 `dist` 目录的内容，而不是 `src`、`public` 或项目根目录。

## 7. 检查正式构建结果

构建完成后启动生产预览：

```powershell
Set-Location -LiteralPath 'E:\blog\astro-blog'
npm run preview -- --host 127.0.0.1
```

再次打开：

```text
http://127.0.0.1:4321/
```

正式预览不会像开发服务器一样自动重新构建。修改源码后需要重新执行：

```powershell
npm run build
```

基础产物检查：

```powershell
(Get-ChildItem -LiteralPath '.\dist' -Recurse -Filter '*.html').Count
(Get-ChildItem -LiteralPath '.\dist\media\music' -Recurse -Filter '*.mp3').Count
(Get-ChildItem -LiteralPath '.\dist\media\music' -Recurse -Filter '*.lrc').Count
```

当前基线应为：

- 21 个 HTML 页面
- 94 个 MP3
- 94 个 LRC
- 7 张专辑封面
- 3 张首页动画宣传图

还需要人工检查：

1. 首页 3D 唱片廊的巡航、拖拽、惯性和循环。
2. 左下角播放器的展开、快速收回和“固定 / 已固定”按钮。
3. 日文原文与中文译文是否正确配对。
4. 七个专辑折叠组是否能播放正确歌曲。
5. 播放中切换首页、音乐页和文章页是否保持进度。
6. 375px 手机宽度下是否存在遮挡或横向溢出。
7. 检查“星河”光影的 Summer / Night 两种状态，以及首页、音乐页和文章页的实际效果。
8. 首次访问是否跟随系统亮暗主题；手动切换后刷新和跨页是否保持 `nabunana:environment-v1`。
9. 光影持续观察至少 30 秒，确认只有极慢呼吸和轻微视差，正文对比度不受影响；减少动画设置下应完全静止。

环境光不使用 Three.js 或 WebGL，也不应出现大于 500KB 的光影专用 JavaScript chunk。正式方案固定为“星河”，右下角只保留 Summer / Night 控件，不显示诗句。

可直接检查指定组合：

```text
http://127.0.0.1:4321/?lightingPreset=silver-river&environment=night&lightingPreview=1
http://127.0.0.1:4321/music/?lightingPreset=silver-river&environment=summer&lightingPreview=1
```

## 8. GitHub Actions 自动发布

`main` 分支是唯一生产分支。Pull Request 只执行检查与构建；推送或手动运行
`Build and deploy blog` 工作流时，Actions 会构建静态站、通过 SSH 上传，并在服务器上
原子切换 `/var/www/blog`。GitHub Pages 不参与发布。

先在 GitHub 仓库的 `Settings → Environments` 创建 `production` 环境，再设置以下 Secrets：

| 名称 | 内容 |
| --- | --- |
| `SERVER_HOST` | 服务器地址，例如 `39.108.101.149` |
| `SERVER_USER` | 有权写入站点父目录的 SSH 用户 |
| `SERVER_PORT` | SSH 端口；不设置时使用 `22` |
| `SERVER_SSH_KEY` | 专用于 Actions 的 SSH 私钥全文 |
| `SERVER_KNOWN_HOSTS` | `ssh-keyscan -H <服务器地址>` 的输出；应在可信终端核对指纹 |

服务器准备完成后，在仓库 `Settings → Secrets and variables → Actions → Variables` 中新增
Repository Variable `DEPLOY_ENABLED=true`。没有这个开关时，工作流只构建和测试，不会连接服务器；
这样首次推送 `main` 时不会因为凭据尚未配置而误发布。

可选的 Environment Variables：

| 名称 | 默认值 | 用途 |
| --- | --- | --- |
| `DEPLOY_PATH` | `/var/www/blog` | Nginx 站点根目录 |
| `SITE_URL` | `https://elma-gohan.xyz` | 发布后的 HTTPS 健康检查地址 |

MP3 不进入 Git 仓库。自动发布会从服务器当前版本继承 `media/music` 下的 MP3；如果
服务器完全没有 MP3，发布会停止且不会切换旧站。首次上线前应手动把音乐文件上传到
`/var/www/blog/media/music`。首次自动发布会把原有实体目录移到相邻的
`.blog-backups/legacy-<UTC 时间>`，以后每个版本位于 `.blog-releases/`，站点目录本身
是指向当前版本的符号链接。

如需回滚，在服务器上把 `/var/www/blog` 符号链接切换到上一版本；每个新版本中的
`PREVIOUS_RELEASE` 文件记录了上一个目标。工作流不会自动删除旧版本或首次发布备份。

## 9. 手动发布边界

本地构建成功不等于可以直接上线。当前流程要求：

1. 完成 `npm run build`。
2. 使用 `npm run preview` 人工验收。
3. 明确确认“验收通过”。
4. 在服务器创建 `/var/www/blog` 的完整回滚备份。
5. 上传新的 `dist` 产物并原子替换网站目录。
6. 检查 Nginx、页面、歌词、图片和音频 Range 请求。
7. 记录备份路径，异常时立即回滚。

不要在未备份时直接清空或覆盖 `/var/www/blog`。服务器密码也不要写入脚本、本文档或命令参数。

## 10. 常见问题

### `npm` 或 `node` 找不到

安装 Node.js 22/24 后重新打开 PowerShell，再运行：

```powershell
node --version
npm --version
```

### 依赖损坏或版本不一致

优先重新执行：

```powershell
npm ci
```

不要手动修改 `node_modules`。

### 构建提示文章字段错误

检查文章顶部 `---` 之间的字段是否符合 `src/content.config.ts`。常见问题包括日期格式错误、`tags` 不是数组、项目状态不在允许值内。

### 修改后预览没有变化

开发模式确认使用的是 `npm run dev`。生产预览模式需要先停止服务、重新构建，再启动预览：

```powershell
npm run build
npm run preview -- --host 127.0.0.1
```

### 音乐或歌词返回 404

检查 `src/data/music.ts` 中的 `slug`、曲目顺序和 `public/media/music` 下的文件编号是否完全一致。路径和文件名区分字符与全角/半角形式。

### 端口已被占用

改用其他端口：

```powershell
npm run preview -- --host 127.0.0.1 --port 4322
```
