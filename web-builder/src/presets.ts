import { defaultTheme, Theme } from './theme';

/** Deep-clones defaultTheme() and applies overrides — keeps presets
 * concise while guaranteeing every field is always populated. */
function preset(overrides: (t: Theme) => void): Theme {
  const t = defaultTheme();
  overrides(t);
  return t;
}

export const presets: Record<string, Theme> = {
  cyberpunk: preset((t) => {
    t.meta = { name: 'cyberpunk', version: '1.0.0', author: '', description: 'Neon-soaked terminal theme' };
    t.banner.text = 'SYSTEM ONLINE // {user}';
    t.graphics.effects.banner = 'glitch';
  }),

  ocean: preset((t) => {
    t.meta = { name: 'ocean', version: '1.0.0', author: '', description: 'Calm blue-green waves' };
    t.colors = {
      primary: '#00B4D8',
      secondary: '#90E0EF',
      background: '#03045E',
      foreground: '#CAF0F8',
      accent: '#0077B6',
      error: '#EF476F',
      success: '#06D6A0',
      warning: '#FFD166',
      muted: '#5C7A99',
    };
    t.prompt.symbol = '~';
    t.graphics.gradient = { enabled: true, from: '#00B4D8', to: '#03045E', direction: 'horizontal' };
    t.banner.text = '~ welcome aboard, {user} ~';
    t.graphics.effects.banner = 'none';
  }),

  minimal: preset((t) => {
    t.meta = { name: 'minimal', version: '1.0.0', author: '', description: 'Grayscale, no distractions' };
    t.colors = {
      primary: '#FFFFFF',
      secondary: '#CCCCCC',
      background: '#000000',
      foreground: '#EEEEEE',
      accent: '#999999',
      error: '#FF6B6B',
      success: '#6BCB77',
      warning: '#FFD93D',
      muted: '#555555',
    };
    t.prompt.symbol = '$';
    t.prompt.format = '{dir} {symbol}';
    t.banner.enabled = false;
    t.graphics.gradient.enabled = false;
    t.graphics.effects = { banner: 'none', prompt: 'none' };
    t.borders.style = 'line';
  }),

  forest: preset((t) => {
    t.meta = { name: 'forest', version: '1.0.0', author: '', description: 'Earthy greens and browns' };
    t.colors = {
      primary: '#52B788',
      secondary: '#95D5B2',
      background: '#1B2A20',
      foreground: '#D8F3DC',
      accent: '#74C69D',
      error: '#E76F51',
      success: '#40916C',
      warning: '#F4A261',
      muted: '#4A5D53',
    };
    t.prompt.symbol = '🌿';
    t.graphics.gradient = { enabled: true, from: '#52B788', to: '#1B2A20', direction: 'vertical' };
    t.banner.text = 'rooted in {user}';
  }),
};
