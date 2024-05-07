function middleware(to, from, next) {
  if (sessionStorage.getItem("accessToken")) {
    next();
  } else {
    next("/");
  }
}

const originRoutes = [
  {
    grant: true,
    path: "/",
    name: "DashboardLayout",
    component: () => import("../layouts/full-screen.vue"),
    children: [
      {
        grant: true,
        path: "",
        name: "Dashboard",
        component: () => import("../pages/index.vue"),
      },
    ],
  },
  {
    grant: true,
    path: "/auth",
    name: "FullScreenLayout",
    component: () => import("@/layouts/full-screen.vue"),
    children: [
      {
        grant: true,
        path: "signin",
        name: "SignIn",
        component: () => import("@/pages/auth/signin.vue"),
      },
      {
        grant: true,
        path: "signup",
        name: "SignUp",
        component: () => import("@/pages/auth/signup.vue"),
      },
    ],
  },
  /* {
      grant: true,
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../pages/NotFound.vue')
} */
];

export const filterByGrant = () => {
  const recuresiveMap = (routes = originRoutes) => {
    return routes
      .map((_route) => {
        if (_route.grant) {
          return _route;
        }
      })
      .filter((v) => !!v);
  };

  return recuresiveMap();
};

export const routes = () => filterByGrant();
