/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#fed7aa",
        "primary-dark": "#facc15",
        secondary: "#0369a1",
        "secondary-dark": "#0c4a6e",
      },
    },
  },
  plugins: [],
};
