import axios from "axios";

export const useCompanyStore = defineStore("company-store", {
  state: () => ({
    companyList: [],
    companyData: null,
    //* company data
    companyId: null,
    companyName: "No Data",
    companyLogoUrl: "defaultLogo.png",
    companyEmail: "No Data",
    companyAddress: "No Data",
    companyPhone: "No Data",
    companyWebsite: "No Data",
    //* related data
    companyDepartments: [],
    companyTitles: [],
    companyLocations: [],
    //* store
    isLoading: false,
    lastFetch: null,
  }),
  getters: {
    //* all data
    getCompanyList: (state) => state.companyList,
    getCompanyData: (state) => state.companyData,
    //* company data
    getCompanyId: (state) => state.companyData?.ID,
    getCompanyName: (state) => state.companyData?.name,
    getCompanyLogoUrl: (state) => state.companyData?.logoUrl,
    getCompanyEmail: (state) => state.companyData?.email,
    getFullAddress: (state) =>
      state.companyData?.address +
      ", " +
      state.companyData?.city +
      ", " +
      state.companyData?.state +
      ", " +
      state.companyData?.country +
      " " +
      state.companyData?.postalCode,
    getCompanyAddress: (state) => state.companyData?.address,
    getCompanyCity: (state) => state.companyData?.city,
    getCompanyState: (state) => state.companyData?.state,
    getCompanyCountry: (state) => state.companyData?.country,
    getCompanyPostalCode: (state) => state.companyData?.postalCode,
    getCompanyPhone: (state) => state.companyData?.phone,
    getCompanyWebsite: (state) => state.companyData?.website,
    //* related data
    getCompanyDepartments: (state) => state.companyDepartments,
    getCompanyTitles: (state) => state.companyTitles,
    getCompanyLocations: (state) => state.companyLocations,
  },
  actions: {
    //! Company API Calls
    async fetchDefaultCompanyData() {
      try {
        const response = await axios.get(
          "http://localhost:3001/api/v1/company",
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
    async fetchOneCompanyData(companyId) {
      try {
        const response = await axios.get(
          `http://localhost:3001/api/v1/company/${companyId}`,
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
    async createCompany({
      companyName,
      companyPhone = null,
      companyWebsite = null,
      companyAddress = null,
      companyCity = null,
      companyState = null,
      companyCountry = null,
      companyPostalCode = null,
      companyLogoUrl = "defaultLogo.png",
    }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.post(
          "http://localhost:3001/api/v1/company",
          {
            name: companyName,
            phone: companyPhone,
            website: companyWebsite,
            address: companyAddress,
            city: companyCity,
            state: companyState,
            country: companyCountry,
            postalCode: companyPostalCode,
            logoUrl: companyLogoUrl,
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
    async deleteCompany(companyId) {
      const messageStore = useMessageStore();
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
      }
    },
    async updateCompany({
      companyId,
      companyName,
      companyPhone = null,
      companyWebsite = null,
      companyAddress = null,
      companyCity = null,
      companyState = null,
      companyCountry = null,
      companyPostalCode = null,
      companyLogoUrl = null,
    }) {
      const messageStore = useMessageStore();
      try {
        const response = await axios.put(
          `http://localhost:3001/api/v1/company/${companyId}`,
          {
            name: companyName,
            phone: companyPhone,
            website: companyWebsite,
            address: companyAddress,
            city: companyCity,
            state: companyState,
            country: companyCountry,
            postalCode: companyPostalCode,
            logoUrl: companyLogoUrl,
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
        // Something happened in setting up the request that triggered an Error
        console.error("Error", error.message);
        messageStore.setError("Something went wrong.");
      }
      console.error(error.config);
    },
    handleSuccess(response) {
      if (response.status >= 200 && response.status < 300) {
        this.companyList = response.data.data.companyList;
        this.companyData = response.data.data.expandedCompany;
        this.companyDepartments =
          response.data.data.expandedCompany?.departments;
        this.companyTitles = response.data.data.expandedCompany?.titles;
        this.companyLocations = response.data.data.expandedCompany?.locations;
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
