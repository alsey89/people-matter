function requireAuth(to, from, next) {
  if (sessionStorage.getItem("accessToken")) {
    next();
  } else {
    next("/auth/signin");
  }
}

const routes = [
  {
    path: "/auth",
    name: "AuthLayout",
    component: () => import("@/layouts/full-screen.vue"),
    children: [
      {
        path: "signin",
        name: "SignIn",
        component: () => import("@/pages/auth/signin.vue"),
      },
    ],
  },
];

export default routes;
