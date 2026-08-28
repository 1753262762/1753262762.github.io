import { mkdir, readFile, readdir, writeFile, copyFile } from 'node:fs/promises';
import path from 'node:path';
import subsetFont from 'subset-font';

const root = process.cwd();
const output = path.join(root, 'public', 'fonts');
const textRoots = ['src', path.join('public', 'media', 'music')];
const extensions = new Set(['.astro', '.css', '.js', '.ts', '.md', '.lrc']);

async function collectFiles(directory, files = []) {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const full = path.join(directory, entry.name);
    if (entry.isDirectory()) await collectFiles(full, files);
    else if (extensions.has(path.extname(entry.name))) files.push(full);
  }
  return files;
}

const files = (await Promise.all(textRoots.map((name) => collectFiles(path.join(root, name))))).flat();
const corpus = [...new Set((await Promise.all(files.map((file) => readFile(file, 'utf8')))).join(''))].join('');
const source = await readFile(path.join(root, 'node_modules', '@fontsource', 'lxgw-wenkai', 'files', 'lxgw-wenkai-latin-500-normal.woff2'));
const subset = await subsetFont(source, corpus, { targetFormat: 'woff2' });

await mkdir(output, { recursive: true });
await writeFile(path.join(output, 'lxgw-wenkai-500-subset.woff2'), subset);
await copyFile(path.join(root, 'node_modules', '@fontsource', 'lxgw-wenkai', 'LICENSE'), path.join(output, 'OFL-LXGW-WenKai.txt'));
await copyFile(path.join(root, 'node_modules', '@fontsource-variable', 'noto-sans-sc', 'LICENSE'), path.join(output, 'OFL-Noto-Sans-SC.txt'));
console.log(`Generated LXGW WenKai subset: ${(subset.length / 1024).toFixed(1)} KiB, ${corpus.length} glyph inputs`);
