export type EnvironmentMode = 'summer' | 'night';

export type LightingPreset = 'dawn-window' | 'pine-shadow' | 'silver-river';

export type LightingIntensity = 'display' | 'standard' | 'reading';

export interface LightingPresetDefinition {
  id: LightingPreset;
  index: string;
  label: string;
  englishLabel: string;
  description: string;
}

export const ENVIRONMENT_STORAGE_KEY = 'nabunana:environment-v1';
export const DEFAULT_LIGHTING_PRESET: LightingPreset = 'dawn-window';

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
