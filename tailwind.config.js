/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [    './pages/**/*.{html,js}',
  './views/**/*.{html,js}',],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}

