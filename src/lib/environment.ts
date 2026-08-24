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

export const ENVIRONMENT_STORAGE_KEY = 'nabunana:environment-v1';
export const DEFAULT_LIGHTING_PRESET: LightingPreset = 'dawn-window';
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

export const lightingRegisterMix: Record<
  Exclude<LightingRegister, 'lab'>,
  { dawn: number; pine: number; river: number; dawnDrawing: 'full' | 'skeleton' | 'quiet' }
> = {
  threshold: { dawn: 0.70, pine: 0.10, river: 0.20, dawnDrawing: 'full' },
  grove: { dawn: 0.40, pine: 0.60, river: 0, dawnDrawing: 'skeleton' },
  field: { dawn: 0.25, pine: 0, river: 0.70, dawnDrawing: 'skeleton' },
  quiet: { dawn: 1, pine: 0, river: 0, dawnDrawing: 'quiet' },
};
