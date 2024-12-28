import posthog from "posthog-js";

export default {
  install(app) {
    app.config.globalProperties.$posthog = posthog.init(
      "phc_Lpa0eynYebYi4OZfer8q1FlERmHtsPfcf72Qp60Ufb0",
      {
        api_host: "https://us.i.posthog.com",
        capture_pageview: false,
      }
    );
  },
};
