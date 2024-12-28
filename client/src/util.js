const redirectAfterDelay = (router, url = "/", delay = 1000) => {
  setTimeout(() => {
    router.push(url);
  }, delay);
};

export { redirectAfterDelay };
