import actions from "./actions";
import { defineStore } from "pinia";
import axios from "axios";

const api_url = import.meta.env.API_URL;

export const useUserStore = defineStore("user-store", {
  state: () => ({
    error_handle: null,
    status: null,
    total: 0,
    item: "",
    items: [],
  }),
  getters: {
    getErrorHandle: (state) => state.error_handle,
    getStatus: (state) => state.status,
    getTotal: (state) => state.total,
    getItem: (state) => state.item,
    getItems: (state) => state.items,
  },
  actions: {
    async signup(payload) {
      try {
        let data = {
          name: payload.name,
          email: payload.email,
          password: payload.password,
          confirmPassword: payload.confirmPassword,
        };
        const response = await axios.post(api_url + "/auth/signup", data);
        if (response.status >= 200 && response.status < 300) {
          if (response.data) {
            sessionStorage.setItem("id", response.data.uid);
            sessionStorage.setItem("email", response.data.email);
            sessionStorage.setItem("name", response.data.name);
            return response.data;
          }
        }
      } catch (err) {
        return err;
      }
    },
    async signin(payload) {
      try {
        let data = {
          email: payload.email,
          password: payload.password,
        };
        const response = await axios.post(api_url + "/auth/signin", data);
        if (response.status >= 200 && response.status < 300) {
          if (response.data) {
            sessionStorage.setItem("id", response.data.uid);
            sessionStorage.setItem("email", response.data.email);
            sessionStorage.setItem("name", response.data.name);
            return response.data;
          }
        }
      } catch (err) {
        return err;
      }
    },
    async getCsrfToken() {
      try {
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
