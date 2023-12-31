// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devTools: true,
  app: {
    head: {
      link: [{ rel: "icon", type: "image/png", href: "/logo.png" }],
    },
  },
  runtimeConfig: {
    public: {
      apiUrl: "http://localhost:3001/api/v1",
    },
  },
  ssr: false,
  devtools: { enabled: true },
  modules: [
    "@nuxtjs/tailwindcss",
    "@pinia/nuxt",
    "nuxt-icon",
    "@pinia-plugin-persistedstate/nuxt",
    "@formkit/auto-animate/nuxt",
  ],
});
