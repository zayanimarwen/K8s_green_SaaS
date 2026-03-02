import type { Config } from 'tailwindcss'

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#f0f9f0', 100: '#dcf0dc', 200: '#b8e0b8',
          300: '#84c984', 400: '#4aab4a', 500: '#228b22',
          600: '#1a6e1a', 700: '#145514', 800: '#0f400f', 900: '#0a2e0a',
        },
        navy: {
          50: '#eef2ff', 500: '#0f3460', 600: '#0d2b52', 800: '#0a1f3d', 900: '#06142a',
        }
      }
    }
  },
  plugins: []
} satisfies Config
