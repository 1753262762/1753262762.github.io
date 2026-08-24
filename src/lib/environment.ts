export type EnvironmentMode = 'summer' | 'night';

export type LightingPreset = 'dawn-window' | 'pine-shadow' | 'silver-river';

export type LightingIntensity = 'display' | 'standard' | 'reading';

export type LightingRegister = 'threshold' | 'grove' | 'field' | 'quiet' | 'lab';

export interface LightingPresetDefinition {
  id: LightingPreset;
  index: string;
  label: string;
  englishLabel: string;
  description: string;
}

export interface LightingLiteraryLine {
  line: string;
  work: string;
  author: string;
  href: string;
}

export const ENVIRONMENT_STORAGE_KEY = 'nabunana:environment-v1';
export const DEFAULT_LIGHTING_PRESET: LightingPreset = 'silver-river';
export const DEFAULT_LIGHTING_REGISTER: LightingRegister = 'quiet';

export const lightingPresets: LightingPresetDefinition[] = [
  {
    id: 'dawn-window',
    index: '01',
    label: '朝窗',
    englishLabel: 'Dawn Window',
    description: '低角度晨光与清水般铺开的月色。',
  },
  {
    id: 'pine-shadow',
    index: '02',
    label: '松影',
    englishLabel: 'Pine Shadow',
    description: '木漏日的光斑与松枝间的银色月光。',
  },
  {
    id: 'silver-river',
    index: '03',
    label: '星河',
    englishLabel: 'Silver River',
    description: '淡紫晨雾与乳白银河中的微光。',
  },
];

export const lightingPresetIds = lightingPresets.map((preset) => preset.id);

export const lightingRegisterIds: LightingRegister[] = ['threshold', 'grove', 'field', 'quiet', 'lab'];

export const isRegister = (value: string | null | undefined): value is LightingRegister =>
  lightingRegisterIds.includes(value as LightingRegister);

export const lightingLiterature: Record<LightingPreset, Record<EnvironmentMode, LightingLiteraryLine>> = {
  'dawn-window': {
    summer: { line: '大きな赤い日', work: '夢十夜', author: '夏目漱石', href: 'https://www.aozora.gr.jp/cards/000148/files/799_14972.html' },
    night: { line: '皎皎空中孤月轮', work: '春江花月夜', author: '张若虚', href: 'https://zh.wikisource.org/wiki/%E6%98%A5%E6%B1%9F%E8%8A%B1%E6%9C%88%E5%A4%9C_(%E5%BC%B5%E8%8B%A5%E8%99%9B)' },
  },
  'pine-shadow': {
    summer: { line: '迟日江山丽', work: '绝句二首', author: '杜甫', href: 'https://zh.wikisource.org/wiki/%E7%B5%95%E5%8F%A5%E4%BA%8C%E9%A6%96_(%E6%9D%9C%E7%94%AB)' },
    night: { line: '明月松间照', work: '山居秋暝', author: '王维', href: 'https://zh.wikisource.org/wiki/%E5%B1%B1%E5%B1%85%E7%A7%8B%E6%9A%9F' },
  },
  'silver-river': {
    summer: { line: '春はあけぼの', work: '枕草子', author: '清少納言', href: 'https://ja.wikisource.org/wiki/%E6%9E%95%E8%8D%89%E5%AD%90_(Wikisource)/%E7%AC%AC%E4%B8%80%E6%AE%B5' },
    night: { line: 'ぼんやりと白い銀河', work: '銀河鉄道の夜', author: '宮沢賢治', href: 'https://www.aozora.gr.jp/cards/000081/files/43737_19215.html' },
  },
};
