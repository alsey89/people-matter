// function middleware(to, from, next) {
//   if (sessionStorage.getItem("accessToken")) {
//     next();
//   } else {
//     next("/");
//   }
// }

const originRoutes = [
  {
    isActive: true,
    path: "/",
    name: "DashboardLayout",
    component: () => import("../layouts/full-screen.vue"),
    children: [
      {
        isActive: true,
        path: "",
        name: "Dashboard",
        component: () => import("../pages/index.vue"),
      },
    ],
  },
  {
    isActive: true,
    path: "/onboarding",
    name: "OnboardingLayout",
    component: () => import("../layouts/full-screen.vue"),
    children: [
      {
        isActive: true,
        path: "",
        name: "Onboarding",
        component: () => import("../pages/onboarding/index.vue"),
      },
      {
        isActive: true,
        path: "confirmation",
        name: "Confirmation",
        component: () => import("../pages/onboarding/confirmation.vue"),
      },
    ],
  },
  {
    isActive: true,
    path: "/auth",
    name: "FullScreenLayout",
    component: () => import("@/layouts/full-screen.vue"),
    children: [
      {
        isActive: true,
        path: "signin",
        name: "SignIn",
        component: () => import("@/pages/auth/signin.vue"),
      },
      // {
      //   isActive: false,
      //   path: "signup",
      //   name: "SignUp",
      //   component: () => import("@/pages/auth/signup.vue"),
      // },
    ],
  },
  /* {
      isActive: true,
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('../pages/NotFound.vue')
} */
];

export const filterByActive = () => {
  const recuresiveMap = (routes = originRoutes) => {
    return routes
      .map((_route) => {
        if (_route.isActive) {
          return _route;
        }
      })
      .filter((v) => !!v);
  };

  return recuresiveMap();
};

export const routes = () => filterByActive();
