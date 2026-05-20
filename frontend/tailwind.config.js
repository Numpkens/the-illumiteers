/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{svelte,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        parchment: {
          DEFAULT: '#F4EFE6',
          light: '#FAF5E6',
        },
        inkpurple: {
          DEFAULT: '#4A3565',
          dark: '#35254A',
          light: '#7C5EA5',
        },
        pastel: {
          pink: '#FBCFE8',
          blue: '#BAE6FD',
          violet: '#DDD6FE',
          amber: '#FDE68A',
        }
      },
      aspectRatio: {
        'card': '2.5 / 3.5',
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
