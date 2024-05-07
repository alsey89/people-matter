import { defineStore } from "pinia";
import axios from "axios";

const api_url = import.meta.env.VITE_API_URL;

export const useUserStore = defineStore("user-store", {
  state: () => ({
    error: "",
    status: "",
    total: 0,
    item: "",
    items: [],
  }),
  getters: {
    getError: (state) => state.error,
    getStatus: (state) => state.status,
    getTotal: (state) => state.total,
    getItem: (state) => state.item,
    getItems: (state) => state.items,
  },
  actions: {
    async createCompany(payload, router) {
      this.error = "";

      const data = {
        companyName: payload.companyName,
        companySize: payload.companySize,
        adminEmail: payload.adminEmail,
        password: payload.password,
        confirmPassword: payload.confirmPassword,
      };
      axios.defaults.headers.post["Content-Type"] = "application/json";
      axios.defaults.headers.post["Accept"] = "application/json";
      axios.defaults.withCredentials = true;

      try {
        const response = await axios.post(api_url + "/company", data);
        if (response.status >= 200 && response.status < 300) {
          return true;
        } else {
          return false;
        }
      } catch (error) {
        if (error.response) {
          if (error.response.status === 409) {
            this.error = "User already exists. Please login.";
            setTimeout(() => {
              router.push("/auth/signin");
            }, 3000);
          } else {
            this.error =
              "Something went wrong. Please try again or contact support.";
          }
        }
        throw error;
      }
    },
    async signin(payload) {
      this.error = "";

      let data = {
        email: payload.email,
        password: payload.password,
      };
      axios.defaults.headers.post["Content-Type"] = "application/json";
      axios.defaults.headers.post["Accept"] = "application/json";
      axios.defaults.withCredentials = true;
      try {
        const response = await axios.post(api_url + "/auth/signin", data);
        if (response.status >= 200 && response.status < 300) {
          return true;
        } else {
          return false;
        }
      } catch (error) {
        if (error.response) {
          if (error.response.status === 401 || error.response.status === 404) {
            this.error = "Invalid email or password. Please try again.";
          } else {
            this.error =
              "Something went wrong. Please try again or contact support.";
          }
        }
        throw error;
      }
    },
    async getCsrfToken() {
      try {
        axios.defaults.withCredentials = true;
        const response = await axios.get(api_url + "/auth/csrf");
        if (response.status >= 200 && response.status < 300) {
          return true;
        }
      } catch (err) {
        return err;
      }
    },
  },
});
