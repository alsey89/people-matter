import { createApp } from "vue";
import "./assets/index.css";
import App from "./App.vue";

import { createPinia } from "pinia";
import initRouter from "@/router";

(async () => {
  const app = createApp(App);

  const router = initRouter();
  app.use(router);

  const pinia = createPinia();
  app.use(pinia);

  app.mount("#app");
})();
