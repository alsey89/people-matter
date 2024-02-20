import axios from "axios";

export const useCompanyStore = defineStore("company-store", {
  state: () => ({
    // companyList: [],
    companyData: null,
    //* company data
    companyId: null,
    companyName: null,
    companyLogoUrl: null,
    companyEmail: null,
    companyAddress: null,
    companyCity: null,
    companyState: null,
    companyCountry: null,
    companyPostalCode: null,
    companyPhone: null,
    companyWebsite: null,
    //* related data
    companyDepartments: null,
    companyLocations: null,
    companyJobs: null,
    //* store
    isLoading: true,
    lastFetch: null,
  }),
  getters: {
    //* all data
    getCompanyData: (state) => state.companyData,
    //* company data
    getCompanyName: (state) => state.companyName, //persisted
    getCompanyLogoUrl: (state) => state.companyLogoUrl, //persisted
    getCompanyId: (state) => state.companyId,
    getCompanyEmail: (state) => state.companyEmail,
    getFullAddress: (state) => {
      let formattedAddress = "";
      let parts = [
        state.companyAddress,
        state.companyCity,
        state.get,
        state.companyCountry,
        state.companyPostalCode,
      ];
      parts.forEach((part) => {
        if (part) {
          formattedAddress += part + " ";
        }
      });
      return formattedAddress;
    },
    getCompanyAddress: (state) => state.companyAddress,
    getCompanyCity: (state) => state.companyCity,
    getCompanyState: (state) => state.companyState,
    getCompanyCountry: (state) => state.companyCountry,
    getCompanyPostalCode: (state) => state.companyPostalCode,
    getCompanyPhone: (state) => state.companyPhone,
    getCompanyWebsite: (state) => state.companyWebsite,
    //* related data
    getCompanyDepartments: (state) => state.companyDepartments,
    getCompanyLocations: (state) => state.companyLocations,
    getCompanyJobs: (state) => state.companyJobs,
    //* get by id
    getDepartmentNameById: (state) => (departmentId) =>
      state.companyDepartments.find(
        (department) => department.ID === departmentId
      ).name,
    getLocationNameById: (state) => (locationId) =>
      state.companyLocations.find((location) => location.ID === locationId)
        .name,
    getJobTitleById: (state) => (jobId) =>
      state.companyJobs.find((job) => job.ID === jobId).title,
    getManagerJobById: (state) => (managerId) =>
      state.companyJobs.find((job) => job.ID === managerId)?.title,
    //* store
    getIsLoading: (state) => state.isLoading,
  },
  actions: {
    async clearCompanyStore() {
      this.state = this.$reset();
    },
    //! Company API Calls
    async fetchCompany() {
      this.isLoading = true;
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/company/default",
          {
            headers: {
              "Content-Type": "application/json",
              Accept: "application/json",
            },
            withCredentials: true,
          }
        );
        const isSuccess = this.handleSuccess(response);
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    async createCompany({ companyFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/company",
          {
            name: companyFormData.companyName,
            phone: companyFormData.companyPhone,
            email: companyFormData.companyEmail,
            website: companyFormData.companyWebsite,
            address: companyFormData.companyAddress,
            city: companyFormData.companyCity,
            state: companyFormData.companyState,
            country: companyFormData.companyCountry,
            postalCode: companyFormData.companyPostalCode,
            logoUrl: companyFormData.companyLogoUrl,
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
          messageStore.setMessage("Company created.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async updateCompany({ companyFormData }) {
      const messageStore = useMessageStore();
      const companyId = companyFormData.companyId;
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/company/${companyId}`,
          {
            name: companyFormData.companyName,
            phone: companyFormData.companyPhone,
            email: companyFormData.companyEmail,
            website: companyFormData.companyWebsite,
            address: companyFormData.companyAddress,
            city: companyFormData.companyCity,
            state: companyFormData.companyState,
            country: companyFormData.companyCountry,
            postalCode: companyFormData.companyPostalCode,
            logoUrl: companyFormData.companyLogoUrl,
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
          messageStore.setMessage("Company updated.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async deleteCompany({ companyId }) {
      const messageStore = useMessageStore();
      this.isLoading = true;
      try {
        const response = await axios.delete(
          `http://localhost:3001/api/v1/company/${companyId}`,
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
          messageStore.setMessage("Company deleted.");
        }
      } catch (error) {
        this.handleError(error);
      } finally {
        this.isLoading = false;
      }
    },
    //! Department API Calls
    async createDepartment({ companyId, departmentFormData }) {
      const messageStore = useMessageStore();
      console.log("department form data: ", departmentFormData);
      try {
        const response = await axios.post(
          `http://localhost:3001/api/v1/company/${companyId}/department`,
          {
            name: departmentFormData.departmentName,
            description: departmentFormData.departmentDescription,
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
          messageStore.setMessage("Department created.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async updateDepartment({ companyId, departmentFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/company/${companyId}/department/${departmentFormData.departmentId}`,
          {
            name: departmentFormData.departmentName,
            description: departmentFormData.departmentDescription,
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
          messageStore.setMessage("Department updated.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async deleteDepartment({ companyId, departmentId }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.delete(
          `http://localhost:3001/api/v1/company/${companyId}/department/${departmentId}`,
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
          messageStore.setMessage("Department deleted.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    //! Location API Calls
    async createLocation({ companyId, locationFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.post(
          `http://localhost:3001/api/v1/company/${companyId}/location`,
          {
            name: locationFormData.locationName,
            phone: locationFormData.locationPhone,
            isHeadOffice: locationFormData.locationIsHeadOffice,
            description: locationFormData.locationDescription,
            address: locationFormData.locationAddress,
            city: locationFormData.locationCity,
            state: locationFormData.locationState,
            country: locationFormData.locationCountry,
            postalCode: locationFormData.locationPostalCode,
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
          messageStore.setMessage("Location created.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async updateLocation({ companyId, locationFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/company/${companyId}/location/${locationFormData.locationId}`,
          {
            name: locationFormData.locationName,
            phone: locationFormData.locationPhone,
            isHeadOffice: locationFormData.locationIsHeadOffice,
            description: locationFormData.locationDescription,
            address: locationFormData.locationAddress,
            city: locationFormData.locationCity,
            state: locationFormData.locationState,
            country: locationFormData.locationCountry,
            postalCode: locationFormData.locationPostalCode,
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
          messageStore.setMessage("Location updated.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async deleteLocation({ companyId, locationId }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.delete(
          `http://localhost:3001/api/v1/company/${companyId}/location/${locationId}`,
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
          messageStore.setMessage("Location deleted.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    //! Job API Calls
    async fetchJobList({ companyId }) {
      try {
        const response = await axios.get(
          `http://localhost:3001/api/v1/job/company/${companyId}`,
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
      }
    },
    async createJob({ companyId, jobFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.post(
          `http://localhost:3001/api/v1/company/${companyId}/job`,
          {
            title: jobFormData.jobTitle,
            departmentId: jobFormData.jobDepartmentId,
            locationId: jobFormData.jobLocationId,
            managerId: jobFormData.jobManagerId,
            MinSalary: jobFormData.jobMinSalary,
            MaxSalary: jobFormData.jobMaxSalary,
            description: jobFormData.jobDescription,
            duties: jobFormData.jobDuties,
            qualifications: jobFormData.jobQualifications,
            experience: jobFormData.jobExperience,
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
          messageStore.setMessage("Job created.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async updateJob({ companyId, jobId, jobFormData }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/company/${companyId}/job/${jobId}`,
          {
            title: jobFormData.jobTitle,
            departmentId: jobFormData.jobDepartmentId,
            locationId: jobFormData.jobLocationId,
            managerId: jobFormData.jobManagerId,
            MinSalary: jobFormData.jobMinSalary,
            MaxSalary: jobFormData.jobMaxSalary,
            description: jobFormData.jobDescription,
            duties: jobFormData.jobDuties,
            qualifications: jobFormData.jobQualifications,
            experience: jobFormData.jobExperience,
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
          messageStore.setMessage("Job updated.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    async deleteJob({ companyId, jobId }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.delete(
          `http://localhost:3001/api/v1/company/${companyId}/job/${jobId}`,
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
          messageStore.setMessage("Job deleted.");
        }
      } catch (error) {
        this.handleError(error);
      }
    },
    //! Common
    handleSuccess(response) {
      const messageStore = useMessageStore();
      if (response.status === 204) {
        this.state = this.$reset();
        messageStore.setMessage("No content.");
        return false;
      }
      if (response.status >= 200 && response.status < 300) {
        this.companyData = response.data.data.expandedCompany;
        this.companyId = response.data.data.expandedCompany.ID;
        this.companyName = response.data.data.expandedCompany.name;
        this.companyEmail = response.data.data.expandedCompany.email;
        this.companyAddress = response.data.data.expandedCompany.address;
        this.companyCity = response.data.data.expandedCompany.city;
        this.companyState = response.data.data.expandedCompany.state;
        this.companyCountry = response.data.data.expandedCompany.country;
        this.companyPostalCode = response.data.data.expandedCompany.postalCode;
        this.companyPhone = response.data.data.expandedCompany.phone;
        this.companyWebsite = response.data.data.expandedCompany.website;
        this.companyLogoUrl = response.data.data.expandedCompany.logoUrl;
        //related data
        this.companyDepartments =
          response.data.data.expandedCompany?.departments;
        this.companyLocations = response.data.data.expandedCompany?.locations;
        this.companyJobs = response.data.data.expandedCompany?.jobs;
        // session storage
        persistedState.sessionStorage.setItem("companyId", this.companyId);
        persistedState.sessionStorage.setItem("companyName", this.companyName);
        persistedState.sessionStorage.setItem(
          "companyLogoUrl",
          this.companyLogoUrl
        );
        return true;
      } else {
        return false;
      }
    },
    handleError(error) {
      const messageStore = useMessageStore();

      if (error.response) {
        console.log(error.response.data);
        switch (error.response.status) {
          case 204:
            messageStore.setError("No content.");
            break;
          case 401:
            messageStore.setError("Invalid credentials.");
            return navigateTo("/signin");
          case 403:
            messageStore.setError("Access denied.");
            return navigateTo("/");
          case 404:
            messageStore.setError("No relevant data.");
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
        console.error(error.request);
        messageStore.setError("No response was received.");
      } else {
        console.log("Error", error.message);
        // Something happened in setting up the request that triggered an Error
        console.error("Error", error.message);
        messageStore.setError("Something went wrong.");
      }
      console.error(error.config);
    },
  },
});
