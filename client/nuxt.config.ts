// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // app: {
  //   head: {
  //     link: [{ rel: "icon", type: "image/png", href: "./assets/logo.png" }],
  //   },
  // },
  runtimeConfig: {
    public: {
      apiUrl: "http://localhost:3001/api/v1",
    },
  },
  devtools: { enabled: true },
  modules: [
    "@nuxtjs/tailwindcss",
    "@pinia/nuxt",
    "nuxt-icon",
    "@pinia-plugin-persistedstate/nuxt",
  ],
});
