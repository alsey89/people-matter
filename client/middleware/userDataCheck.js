export default defineNuxtRouteMiddleware(async (to, from) => {
  if (process.server) return;

  const userStore = useUserStore();
  if (!userStore.getUserData || userStore.shouldFetchUserData()) {
    await userStore.fetchCurrentUserData();
  }
});
