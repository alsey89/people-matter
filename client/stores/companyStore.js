import axios from "axios";

export const useCompanyStore = defineStore("company-store", {
  state: () => ({
    companyList: [],
    companyData: null,
    //* company data
    companyName: null,
    companyLogoUrl: null,
    companyEmail: null,
    companyAddress: null,
    companyPhone: null,
    companyWebsite: null,
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
    getCompanyPhone: (state) => state.companyData?.phone,
    getCompanyWebsite: (state) => state.companyData?.website,
    //* related data
    getCompanyDepartments: (state) => state.companyDepartments,
    getCompanyTitles: (state) => state.companyTitles,
    getCompanyLocations: (state) => state.companyLocations,
  },
  actions: {
    //! Company API Calls
    async fetchCompanyData() {
      const messageStore = useMessageStore();
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
        if (response.status === 200) {
          // messageStore.setMessage("Successfully fetched company data.");
          this.companyList = response.data.data.companyList;
          this.companyData = response.data.data.expandedCompany;
          this.companyDepartments =
            response.data.data.expandedCompany?.departments;
          this.companyTitles = response.data.data.expandedCompany?.titles;
          this.companyLocations = response.data.data.expandedCompany?.locations;
        }
      } catch (error) {
        this.handleError(error);
      }
    },
  },
  // persist: {
  //   storage: persistedState.sessionStorage,
  // },
});
