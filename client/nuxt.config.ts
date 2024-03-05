// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  app: {
    head: {
      link: [
        {
          rel: "icon",
          type: "image/x-icon",
          href: process.env.FAVICON_PATH || "/favicon.ico",
        },
      ],
    },
  },
  runtimeConfig: {
    public: {
      // localhost3001 in local development or the production API URL
      apiBaseUrl: process.env.API_BASE_URL,
    },
  },
  devtools: {
    enabled: true,
    timeline: {
      enabled: true,
    },
  },
  modules: [
    "@nuxtjs/tailwindcss",
    "@pinia/nuxt",
    "nuxt-icon",
    "@pinia-plugin-persistedstate/nuxt",
    "@formkit/auto-animate/nuxt",
    "shadcn-nuxt",
  ],
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: "",
    /**
     * Directory that the component lives in.
     * @default "./components/ui"
     */
    componentDir: "./components/ui",
  },
});
