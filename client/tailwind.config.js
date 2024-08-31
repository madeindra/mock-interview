/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'dark-bg': '#1A1A1A',
        'dark-surface': '#2C2C2C',
        'dark-primary': '#4A90E2',
        'dark-secondary': '#718096',
        'dark-error': '#E53E3E',
        'dark-on-bg': '#E1E1E1',
        'dark-on-surface': '#FFFFFF',
      },
    },
  },
  plugins: [],
}