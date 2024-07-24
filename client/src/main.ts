import { createApp, nextTick } from "vue";
import "./assets/index.css";
import App from "./App.vue";

import { createPinia } from "pinia";
import initRouter from "@/router";
import { autoAnimatePlugin } from "@formkit/auto-animate/vue";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { ScrollToPlugin } from "gsap/ScrollToPlugin";

// import posthogPlugin from "./plugins/posthog";

gsap.registerPlugin(ScrollTrigger, ScrollToPlugin);

(async () => {
  const app = createApp(App);

  const router = initRouter();
  app.use(router);

  const pinia = createPinia();
  app.use(pinia);

  app.use(autoAnimatePlugin);

  //---- Posthog -----
  // app.use(posthogPlugin);
  // router.beforeEach((to, from, next) => {
  //   //capture pageleave
  //   if (from.fullPath !== to.fullPath) {
  //     app.config.globalProperties.$posthog.capture("$pageleave", {
  //       path: from.fullPath,
  //     });
  //   }
  //   next();
  // });
  // router.afterEach((to, from, failure) => {
  //   if (!failure && to.fullPath !== from.fullPath) {
  //     nextTick(() => {
  //       app.config.globalProperties.$posthog.capture("$pageview", {
  //         path: to.fullPath,
  //       });
  //     });
  //   }
  // });

  app.mount("#app");
})();
