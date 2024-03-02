import axios from "axios";

export const useUserStore = defineStore("user-store", {
  state: () => ({
    //* complete data
    currentUserData: null,
    allUsersData: null,
    singleUserData: null,
    //* store
    isLoading: false,
  }),
  getters: {
    //* complete data
    getCurrentUserData: (state) => state.currentUserData,
    getAllUsersData: (state) => state.allUsersData,
    getSingleUserData: (state) => state.singleUserData,
    //* current user data
    getCurrentUserAvatarUrl: (state) => state.currentUserData?.avatarUrl,
    getCurrentUserFullName: (state) =>
      state.currentUserData?.firstName + " " + state.currentUserData?.lastName,
    getCurrentUserUserId: (state) => state.currentUserData?.userId,
    getCurrentUserEmail: (state) => state.currentUserData?.email,
    getCurrentUserIsAdmin: (state) => state.currentUserData?.isAdmin,
    getCurrentUserTitle: (state) =>
      state.currentUserData?.assignedJob?.job?.title?.name,
    getCurrentUserDepartment: (state) =>
      state.currentUserData?.assignedJob?.job?.department?.name,
    //* single user data
    //main
    getSingleUserAvatarUrl: (state) => state.singleUserData?.avatarUrl,
    getSingleUserFullName: (state) => {
      let fullName = "";
      if (state.singleUserData?.firstName) {
        fullName += state.singleUserData?.firstName;
      }
      if (state.singleUserData?.middleName) {
        fullName += " " + state.singleUserData?.middleName;
      }
      if (state.singleUserData?.lastName) {
        fullName += " " + state.singleUserData?.lastName;
      }
      return fullName;
    },
    getSingleUserUserId: (state) => state.singleUserData?.userId,
    getSingleUserEmail: (state) => state.singleUserData?.email,
    getSingleUserRole: (state) => state.singleUserData?.role,
    getSingleUserTitle: (state) =>
      state.singleUserData?.assignedJob?.job?.title?.name,
    getSingleUserDepartment: (state) =>
      state.singleUserData?.assignedJob?.job?.department?.name,
    //contact info
    getSingleUserContactInfo: (state) => state.singleUserData?.contactInfo,
    getSingleUserMobile: (state) => state.singleUserData?.contactInfo?.mobile,
    getSingleUserAddress: (state) => state.singleUserData?.contactInfo?.address,
    getSingleUserCity: (state) => state.singleUserData?.contactInfo?.city,
    getSingleUserState: (state) => state.singleUserData?.contactInfo?.state,
    getSingleUserPostalCode: (state) =>
      state.singleUserData?.contactInfo?.postalCode,
    getSingleUserCountry: (state) => state.singleUserData?.contactInfo?.country,
    getSingleUserFullAddress: (state) => {
      let fullAddress = "";
      if (state.singleUserData?.contactInfo?.address) {
        fullAddress += state.singleUserData?.contactInfo?.address;
      }
      if (state.singleUserData?.contactInfo?.city) {
        fullAddress += ", " + state.singleUserData?.contactInfo?.city;
      }
      if (state.singleUserData?.contactInfo?.state) {
        fullAddress += ", " + state.singleUserData?.contactInfo?.state;
      }
      if (state.singleUserData?.contactInfo?.postalCode) {
        fullAddress += ", " + state.singleUserData?.contactInfo?.postalCode;
      }
      return fullAddress;
    },
    //emergency contact
    getSingleUserEmergencyContact: (state) =>
      state.singleUserData?.emergencyContact,
    getSingleUserEmergencyContactName: (state) =>
      state.singleUserData?.emergencyContact?.firstName,
    getSingleUserEmergencyContactLastName: (state) =>
      state.singleUserData?.emergencyContact?.lastName,
    getSingleUserEmergencyContactMiddleName: (state) =>
      state.singleUserData?.emergencyContact?.middleName,
    getSingleUserEmergencyContactFullName: (state) => {
      let fullName = "";
      if (state.singleUserData?.emergencyContact?.firstName) {
        fullName += state.singleUserData?.emergencyContact?.firstName;
      }
      if (state.singleUserData?.emergencyContact?.middleName) {
        fullName += " " + state.singleUserData?.emergencyContact?.middleName;
      }
      if (state.singleUserData?.emergencyContact?.lastName) {
        fullName += " " + state.singleUserData?.emergencyContact?.lastName;
      }
      return fullName;
    },
    getSingleUserEmergencyContactMobile: (state) =>
      state.singleUserData?.emergencyContact?.mobile,
    getSingleUserEmergencyContactEmail: (state) =>
      state.singleUserData?.emergencyContact?.email,
    getSingleUserEmergencyContactRelation: (state) =>
      state.singleUserData?.emergencyContact?.relation,
    //* store
    getIsLoading: (state) => state.isLoading,
  },
  actions: {
    //! Auth API Calls
    async signin({ email, password }) {
      const runtimeConfig = useRuntimeConfig();
      const apiBaseUrl = runtimeConfig.public.apiBaseUrl;
      const messageStore = useMessageStore();
      const companyStore = useCompanyStore();
      this.isLoading = true;
      try {
        const response = await axios.post(
          `${apiBaseUrl}/api/v1/auth/signin`,
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
        if (response.status === 200) {
          messageStore.setMessage("Successfully signed in.");
          this.currentUserData = response.data.data;
          await companyStore.fetchCompany();
          return navigateTo("/");
        }
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async signup({ username, email, password, confirmPassword }) {
      const messageStore = useMessageStore();
      this.isLoading = true;
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
        if (response.status === 200) {
          messageStore.setMessage("Successfully signed up.");
          this.currentUserData = response.data.data;
          return navigateTo("/");
        }
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async signout() {
      this.isLoading = true;
      try {
        await axios.post(
          "http://localhost:3001/api/v1/auth/signout",
          {},
          {
            withCredentials: true,
          }
        );
        sessionStorage.clear();
        return navigateTo("/auth/signin");
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async getCsrfToken() {
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/auth/csrf",
          {
            withCredentials: true,
          }
        );
      } catch (error) {
        console.error("Error:", error.response);
        this.handleError(error);
      }
    },
    async checkAuth() {
      try {
        await axios.get("http://localhost:3001/api/v1/auth/check", {
          withCredentials: true,
        });
      } catch (error) {
        console.error("Error:", error.response);
        sessionStorage.clear();
        return navigateTo("/auth/signin", { external: true });
      }
    },
    //! User API Calls
    //* current
    async fetchCurrentUserData() {
      this.currentUserData = null;
      this.isLoading = true;
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/user/current",
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        if (response.status === 200) {
          this.currentUserData = response.data.data;
        }
      } catch (error) {
        console.error("Error:", error.response);
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    //* all users
    async fetchAllUsersData() {
      this.allUsersData = null;
      this.isLoading = true;
      try {
        const response = await axios.get(
          `http://localhost:3001/api/v1/user/list`,
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        this.handleSuccess(response);
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async createUser({ userFormData }) {
      const messageStore = useMessageStore();
      this.isLoading = true;
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/user/list",
          {
            email: userFormData.email,
            firstName: userFormData.firstName,
            middleName: userFormData.middleName,
            lastName: userFormData.lastName,
            isAdmin: userFormData.isAdmin,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        const isSuccess = this.handleSuccess(response);
        if (isSuccess) {
          messageStore.setMessage("User created.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async deleteUser({ userId }) {
      const messageStore = useMessageStore();
      this.isLoading = true;
      try {
        const response = await axios.delete(
          `http://localhost:3001/api/v1/user/list/${userId}`,
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        const isSuccess = this.handleSuccess(response);
        if (isSuccess) {
          messageStore.setMessage("User deleted.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async updateUser({ userId, userFormData }) {
      const messageStore = useMessageStore();
      this.isLoading = true;
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/user/list/${userId}`,
          {
            email: userFormData.email,
            firstName: userFormData.firstName,
            middleName: userFormData.middleName,
            lastName: userFormData.lastName,
            isAdmin: userFormData.isAdmin,
          },
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        const isSuccess = this.handleSuccess(response);
        if (isSuccess) {
          messageStore.setMessage("User updated.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    //* single user details
    async fetchSingleUserData(userId) {
      this.singleUserData = null;
      this.isLoading = true;
      try {
        const response = await axios.get(
          `http://localhost:3001/api/v1/user/${userId}`,
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        this.singleUserData = response.data.data;
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    //! Utilities
    handleError(error) {
      console.log("entering handleError");
      const messageStore = useMessageStore();
      if (error.response) {
        console.log(error.response.data);
        switch (error.response.status) {
          case 401:
            messageStore.setError("Invalid credentials.");
            return navigateTo("/auth/signin");
          case 403:
            messageStore.setError("Access denied.");
            return navigateTo("/");
          case 404:
            messageStore.setError("Data not found.");
            break;
          case 409:
            messageStore.setError("Data already exists.");
            break;
          case 500:
            messageStore.setError("Server error.");
            break;
          default:
            messageStore.setError("Something went wrong.");
            return navigateTo("/");
        }
      } else if (error.request) {
        // The request was made but no response was received
        console.log("Request Error", error.request);
        messageStore.setError("No response was received.");
      } else {
        // Something happened in setting up the request that triggered an Error
        console.log("Error", error.message);
        messageStore.setError("Something went wrong.");
      }
      console.log(error.config);
    },
    handleSuccess(response) {
      const messageStore = useMessageStore();
      if (response.status === 204) {
        this.state = this.$reset();
        messageStore.setMessage("No content.");
        return false;
      }
      if (response.status >= 200 && response.status < 300) {
        this.allUsersData = response.data.data;
        return true;
      } else {
        return false;
      }
    },
  },
  // persist: {
  //   storage: persistedState.sessionStorage,
  // },
});
