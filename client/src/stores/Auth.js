import { defineStore } from "pinia";
import { useRouter } from "vue-router";
import axios from "axios";
import api from "@/plugins/axios";
import posthog from "posthog-js";

export const useAuthStore = defineStore("auth-store", {
  state: () => ({
    userRoles: [],
    userId: "",
    email: "",
    error: "",
  }),
  getters: {
    getError: (state) => state.error,
    getUserId: (state) => state.userId,
    getEmail: (state) => state.email,
    isAuthenticated: (state) => state.userId !== "",
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
      try {
        const response = await api.post("/company", data);
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
    async signin(payload, router) {
      this.error = "";
      let data = {
        email: payload.email,
        password: payload.password,
      };
      try {
        const response = await api.post("/auth/signin", data);
        if (response.status >= 200 && response.status < 300) {
          this.userRoles = response.data.data;
          router.push(`select-company?email=${payload.email}`);
        } else {
          return false;
        }
      } catch (error) {
        switch (error.response?.status) {
          case 401:
          case 404:
            this.error = "Invalid email or password. Please try again.";
            break;
          case 403:
            this.error = "Please verify your email address.";
            break;
          default:
            this.error =
              "Something went wrong. Please try again or contact support.";
            break;
        }
        throw error;
      }
    },
    async getJwt(userRoleId, email, router) {
      this.error = "";
      let data = {
        userRoleId: userRoleId,
        email: email,
      };
      try {
        const response = await api.post("/auth/signin/token", data);
        if (response.status >= 200 && response.status < 300) {
          this.userId = response.data.data.id;
          this.email = response.data.data.email;
          posthog.identify(response.data.data.email);
          router.push("/");
        }
      } catch (error) {
        switch (error.response?.status) {
          case 401:
          case 404:
            this.error = "Invalid email or password. Please try again.";
            break;
          case 403:
            this.error = "Please verify your email address.";
            break;
          default:
            this.error =
              "Something went wrong. Please try again or contact support.";
            break;
        }
        throw error;
      }
    },
    async signout() {
      const router = useRouter();
      try {
        await api.post("/auth/signout");
      } catch (error) {
        console.log(error);
      }
      posthog.reset();
      router.push("/auth/signin");
    },
    async getCsrfToken() {
      try {
        const response = await api.get("/auth/csrf");
        if (response.status >= 200 && response.status < 300) {
          return true;
        }
      } catch (err) {
        return err;
      }
    },
  },
});
