import axios from "axios";

export const useUserStore = defineStore("user-store", {
  state: () => ({
    currentUserData: null,
    allUsersData: null,
    userData: null,
    //* user data
    firstName: null,
    middleName: null,
    lastName: null,
    nickName: null,
    //* account data
    email: null,
    avatarUrl: null,
    //* store
    isLoading: false,
    lastFetch: null,
  }),
  getters: {
    //* all user data
    getCurrentUserData: (state) => state.currentUserData,
    getAllUsersData: (state) => state.allUsersData,
    getCurrentUserAvatarUrl: (state) => state.currentUserData?.avatarUrl,
    getCurrentUserFullName: (state) =>
      state.userData?.firstName + " " + state.currentUserData?.lastName,
    getCurrentUserUserId: (state) => state.currentUserData?.userId,
    getCurrentUserEmail: (state) => state.currentUserData?.email,
    getCurrentUserIsAdmin: (state) => state.currentUserData?.isAdmin,
  },
  actions: {
    //! Auth API Calls
    async signin({ store, email, password }) {
      const messageStore = useMessageStore();
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
          messageStore.setMessage("Successfully signed in.");
          return navigateTo("/");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async signup({ store, username, email, password, confirmPassword }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/auth/signup",
          {
            username: username,
            email: email,
            password: password,
            confirmPassword: confirmPassword,
          },
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
        if (response.status === 200) {
          messageStore.setMessage("Successfully signed up.");
          return navigateTo("/");
        }
      } catch (error) {
        this.handleError(error);
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
        this.handleError(error);
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
        const response = await axios.get("http://localhost:3001/api/v1/user", {
          headers: {
            "Content-Type": "application/json",
            Accept: "application/json",
          },
          withCredentials: true,
        });
        this.currentUserData = response.data.data;
        this.firstName = response.data.data.firstName;
        this.lastName = response.data.data.lastName;
        this.middleName = response.data.data.middleName;
        this.nickName = response.data.data.nickName;
        this.email = response.data.data.email;
        this.avatarUrl = response.data.data.avatarUrl;
        this.lastFetch = Date.now();
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async fetchOneUserData(userId) {
      this.isLoading = true;
      try {
        const response = await axios.get(
          "`http://localhost:3001/api/v1/user/${userId}",
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        this.userData = response.data.data;
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async fetchAllUsersData(store) {
      this.isLoading = true;
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/admin/user/all",
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        this.allUsersData = response.data.data;
        this.lastFetch = Date.now();
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    //! Utilities
    shouldFetchUserData() {
      const THRESHOLD = 5 * 60 * 1000; //* 5 minutes
      if (!this.lastFetch) return true;
      return Date.now() - this.lastFetch > THRESHOLD;
    },
    handleError(error) {
      const messageStore = useMessageStore();

      if (error.response) {
        console.log(error.response.data);
        switch (error.response.status) {
          case 401:
            messageStore.setError("Invalid credentials.");
            return navigateTo("/signin");
          case 403:
            messageStore.setError("Access denied.");
            return navigateTo("/");
          case 404:
            messageStore.setError("Data not found.");
          case 409:
            messageStore.setError("Data already exists.");
          case 500:
            messageStore.setError("Server error.");
          default:
            messageStore.setError("Something went wrong.");
            return navigateTo("/");
        }
      } else if (error.request) {
        // The request was made but no response was received
        console.log(error.request);
        messageStore.setError("No response was received.");
      } else {
        // Something happened in setting up the request that triggered an Error
        console.log("Error", error.message);
        messageStore.setError("Something went wrong.");
      }
      console.log(error.config);
    },
  },
  // persist: {
  //   storage: persistedState.sessionStorage,
  // },
});
