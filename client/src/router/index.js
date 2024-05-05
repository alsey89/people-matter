import { createRouter, createWebHistory } from "vue-router";
import routes from "./routes.js";

const initRouter = () => {
  const router = createRouter({
    history: createWebHistory(),
    base: "/",
    routes: routes,
  });

  return router;
};

export default initRouter;
