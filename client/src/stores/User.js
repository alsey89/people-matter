import { defineStore } from "pinia";
import api from "@/plugins/axios";
import { useRouter } from "vue-router";

export const useUserStore = defineStore("user-store", {
  state: () => ({
    users: [],
    selectedUser: null,
    error: "",
  }),
  getters: {
    getSelectedUser: (state) => state.selectedUser,
  },
  actions: {
    //setters
    async setSelectedUser(user) {
      this.selectedUser = user;
    },
    async getUsers() {
      const router = useRouter(); // move useRouter inside the action
      try {
        const response = await api.get("/admin/user");
        this.users = response.data.data;
      } catch (error) {
        switch (error.response?.status) {
          case 400:
            this.error = "Invalid request. Please try again.";
            break;
          case 401:
            this.error = "Unauthorized. Please login.";
            router.push("/auth/signin"); // now router.push will work
            break;
          case 404:
            this.error = "Data not found. Please try again.";
            break;
          default:
            this.error =
              "Something went wrong. Please try again or contact support.";
            break;
        }
        throw error;
      }
    },
  },
});
