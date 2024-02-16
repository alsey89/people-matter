/** @type {import('tailwindcss').Config} */

module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#5c986a",
        "primary-dark": "#528a60",
        "primary-light": "#7acc8f",
        secondary: "#7acc8f",
        "secondary-dark": "#0c4a6e",
        "primary-bg": "#fbf8f3",
        "secondary-bg": "#f4e8d5",
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
};

// module.exports = {
//   content: ["./src/**/*.{js,jsx,ts,tsx}"],
//   theme: {
//     extend: {
//       colors: {
//         primary: "#fed7aa",
//         "primary-dark": "#facc15",
//         secondary: "#0369a1",
//         "secondary-dark": "#0c4a6e",
//         "primary-bg": "#faf9f3",
//         "secondary-bg": "#fffff", //white
//       },
//     },
//   },
//   plugins: [],
// };
