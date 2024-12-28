import { defineStore } from "pinia";
import { useRouter } from "vue-router";
import api from "@/plugins/axios";

export const useCompanyStore = defineStore("company-store", {
  state: () => ({
    company: {},
    logoUrl: "/placeholderImage.png",
    locations: [],
    departments: [],
    positions: [],
    error: "",
  }),
  getters: {
    getError: (state) => state.error,
    getAddress: (state) => {
      if (state.company) {
        return formatAddress(state.company);
      }
      return "";
    },
  },
  actions: {
    async getCompany() {
      this.error = "";
      try {
        const response = await api.get("/company");
        if (response.status >= 200 && response.status < 300) {
          this.company = response.data.data;
          this.locations = response.data.data.locations;
          this.departments = response.data.data.departments;
          this.positions = response.data.data.positions;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async updateCompany(form) {
      this.error = "";
      try {
        const response = await api.put("admin/company", form);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async addLocation(form) {
      this.error = "";
      try {
        const response = await api.post("admin/company/location", form);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async updateLocation(form) {
      this.error = "";
      try {
        const response = await api.put(
          `admin/company/location/${form.id}`,
          form
        );
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async deleteLocation(id) {
      this.error = "";
      try {
        const response = await api.delete(`admin/company/location/${id}`);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async addDepartment(form) {
      this.error = "";
      try {
        const response = await api.post("admin/company/department", form);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async updateDepartment(form) {
      this.error = "";
      try {
        const response = await api.put(
          `admin/company/department/${form.id}`,
          form
        );
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async deleteDepartment(id) {
      this.error = "";
      try {
        const response = await api.delete(`admin/company/department/${id}`);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async addPosition(form) {
      this.error = "";
      try {
        const response = await api.post("admin/company/position", form);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async updatePosition(form) {
      this.error = "";
      try {
        const response = await api.put(
          `admin/company/position/${form.id}`,
          form
        );
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    async deletePosition(id) {
      this.error = "";
      try {
        const response = await api.delete(`admin/company/position/${id}`);
        if (response.status >= 200 && response.status < 300) {
          await this.getCompany();
          return true;
        } else {
          return false;
        }
      } catch (error) {
        this.handleApiError(error);
        throw error;
      }
    },
    handleApiError(error) {
      const router = useRouter();
      if (error.response) {
        switch (error.response.status) {
          case 400:
            this.error = "Invalid request. Please try again.";
            break;
          case 401:
            this.error = "Unauthorized. Please login.";
            router.push("/auth/signin");
            break;
          case 404:
            this.error = "Not found. Please try again.";
            break;
          default:
            this.error =
              "Something went wrong. Please try again or contact support.";
            break;
        }
      } else {
        this.error = "Network error. Please try again or contact support.";
      }
    },
  },
});

//helper functions
const formatAddress = (company) => {
  return [
    company.address,
    company.city,
    company.state,
    company.country,
    company.postalCode,
  ]
    .filter(Boolean)
    .join(", ");
};
