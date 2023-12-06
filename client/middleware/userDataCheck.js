export default defineNuxtRouteMiddleware(async (to, from) => {
  if (process.client) {
    const userStore = useUserStore();
    if (!userStore.getUserData || userStore.shouldFetchUserData()) {
      await userStore.fetchCurrentUserData();
    }
  }
});
