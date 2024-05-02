/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/*.html"],
  safelist: [
    {
      pattern: /from-(red|orange|amber|yellow|green|blue|purple|pink|lime|emerald|teal|cyan|indigo|violet|fuchsia|rose)-\d{3}/,
    },
    {
      pattern: /to-(red|orange|amber|yellow|green|blue|purple|pink|lime|emerald|teal|cyan|indigo|violet|fuchsia|rose)-\d{3}/,
    },
  ],
  theme: {
  },
  plugins: [],
}

