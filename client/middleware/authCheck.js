import axios from "axios";

export default defineNuxtRouteMiddleware(async (to, from) => {
  if (process.server) {
    //* skip serverside, no cookies there
    return;
  }
  if (process.client) {
    const userStore = useUserStore();
    await userStore.checkAuth();
  }
});
