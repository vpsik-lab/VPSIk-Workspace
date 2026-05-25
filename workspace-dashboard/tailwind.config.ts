import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        vpsik: {
          50: '#f0f7ff',
          100: '#e0f0ff',
          200: '#baddff',
          300: '#7cc0ff',
          400: '#36a0ff',
          500: '#0c82f2',
          600: '#0066cc',
          700: '#0052a6',
          800: '#004688',
          900: '#063c71',
        },
      },
    },
  },
  plugins: [],
}

export default config
