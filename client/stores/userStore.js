import axios from "axios";

export const useUserStore = defineStore("user-store", {
  state: () => ({
    userData: null,
    isLoading: false,
    lastFetch: null,
  }),
  getters: {
    //* all user data
    getUserData: (state) => state.userData,
    getAvatarUrl: (state) => state.userData?.avatarUrl,
    getUsername: (state) => state.userData?.username,
    getUserId: (state) => state.userData?.userId,
    getEmail: (state) => state.userData?.email,
    getIsAdmin: (state) => state.userData?.isAdmin,
  },
  actions: {
    //! Auth API Calls
    async signin({ store, email, password }) {
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/auth/signin",
          {
            email: email,
            password: password,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        console.log("Response Data:", response.data);
        this.userData = response.data.data;
        this.lastFetch = Date.now();
        if (response.status === 200) {
          return navigateTo("/");
        }
      } catch (error) {
        console.error("Error:", error);
        return error.response;
      }
    },
    async signup({ store, email, password }) {
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/auth/signup",
          {
            email: email,
            password: password,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        console.log("Response Data:", response.data);
        this.userData = response.data.data;
        this.lastFetch = Date.now();
        if (response.status === 200) {
          return navigateTo("/");
        }
      } catch (error) {
        console.error("Error:", error.response);
        return error.response;
      }
    },
    async signout() {
      try {
        await axios.post(
          "http://localhost:3001/api/v1/auth/signout",
          {},
          {
            withCredentials: true,
          }
        );
        sessionStorage.clear();
        return navigateTo("/signin");
      } catch (error) {
        console.error("Error:", error);
        return error.response;
      }
    },
    async checkAuth(store) {
      try {
        await axios.get("http://localhost:3001/api/v1/auth/check", {
          withCredentials: true,
        });
      } catch (error) {
        console.error("Error:", error.response);
        sessionStorage.clear();
        return navigateTo("/signin", { external: true });
      }
    },
    //! User API Calls
    async fetchCurrentUserData(store) {
      this.isLoading = true;
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/user/data",
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        this.userData = response.data.data;
        this.lastFetch = Date.now();
      } catch (error) {
        console.error("Error:", error.response);
        return error.response;
      } finally {
        this.isLoading = false;
      }
    },
    shouldFetchUserData() {
      const THRESHOLD = 5 * 60 * 1000; //* 5 minutes
      if (!this.lastFetch) return true;
      return Date.now() - this.lastFetch > THRESHOLD;
    },
  },
  persist: {
    storage: persistedState.sessionStorage,
  },
});
