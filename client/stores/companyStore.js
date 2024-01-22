import axios from "axios";

export const useCompanyStore = defineStore("company-store", {
  state: () => ({
    companyList: [],
    companyData: null,
    //* company data
    companyId: null,
    companyName: null,
    companyLogoUrl: "defaultLogo.png",
    companyEmail: null,
    companyAddress: null,
    companyPhone: null,
    companyWebsite: null,
    //* related data
    companyDepartments: null,
    companyLocations: null,
    //* store
    isLoading: false,
    lastFetch: null,
  }),
  getters: {
    //* all data
    getCompanyList: (state) => state.companyList,
    getCompanyData: (state) => state.companyData,
    //* company data
    getCompanyName: (state) => state.companyName, //persisted
    getCompanyLogoUrl: (state) => state.companyLogoUrl, //persisted

    getCompanyId: (state) => state.companyData?.ID,
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
    async fetchCompanyList() {},
    async fetchCompanyListAndExpandDefault() {
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
        this.handleSuccess(response);
      } catch (error) {
        this.handleError(error);
      }
    },
    async fetchCompanyListAndExpandById(companyId) {
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
    //! Title API Calls
    // async createTitle({ companyId, titleFormData }) {
    //   const messageStore = useMessageStore();
    //   try {
    //     const response = await axios.post(
    //       `http://localhost:3001/api/v1/company/${companyId}/title`,
    //       {
    //         name: titleFormData.titleName,
    //         description: titleFormData.titleDescription,
    //         departmentId: titleFormData.departmentId,
    //         departmentName: titleFormData.departmentName,
    //       },
    //       {
    //         headers: {
    //           "Content-Type": "application/json",
    //           Accept: "application/json",
    //         },
    //         withCredentials: true,
    //       }
    //     );
    //     const isSuccess = this.handleSuccess(response);
    //     if (isSuccess) {
    //       messageStore.setMessage("Title created.");
    //     }
    //   } catch (error) {
    //     this.handleError(error);
    //   }
    // },
    // async updateTitle({ companyId, titleFormData }) {
    //   const messageStore = useMessageStore();
    //   try {
    //     const response = await axios.put(
    //       `http://localhost:3001/api/v1/company/${companyId}/title/${titleFormData.titleId}`,
    //       {
    //         name: titleFormData.titleName,
    //         description: titleFormData.titleDescription,
    //         departmentId: titleFormData.departmentId,
    //         departmentName: titleFormData.departmentName,
    //       },
    //       {
    //         headers: {
    //           "Content-Type": "application/json",
    //           Accept: "application/json",
    //         },
    //         withCredentials: true,
    //       }
    //     );
    //     const isSuccess = this.handleSuccess(response);
    //     if (isSuccess) {
    //       messageStore.setMessage("Title updated.");
    //     }
    //   } catch (error) {
    //     this.handleError(error);
    //   }
    // },
    // async deleteTitle({ companyId, titleId }) {
    //   const messageStore = useMessageStore();
    //   try {
    //     const response = await axios.delete(
    //       `http://localhost:3001/api/v1/company/${companyId}/title/${titleId}`,
    //       {
    //         headers: {
    //           "Content-Type": "application/json",
    //           Accept: "application/json",
    //         },
    //         withCredentials: true,
    //       }
    //     );
    //     const isSuccess = this.handleSuccess(response);
    //     if (isSuccess) {
    //       messageStore.setMessage("Title deleted.");
    //     }
    //   } catch (error) {
    //     this.handleError(error);
    //   }
    // },
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
    //! Common
    handleSuccess(response) {
      if (response.status >= 200 && response.status < 300) {
        this.companyList = response.data.data.companyList;
        this.companyData = response.data.data.expandedCompany;
        this.companyId = response.data.data.expandedCompany.ID;
        this.companyName = response.data.data.expandedCompany.name;
        this.companyDepartments =
          response.data.data.expandedCompany?.departments;
        // this.companyTitles = response.data.data.expandedCompany?.titles;
        this.companyLocations = response.data.data.expandedCompany?.locations;
        //store active company in session
        persistedState.sessionStorage.setItem(
          "activeCompanyId",
          +this.companyId
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
        // Something happened in setting up the request that triggered an Error
        console.error("Error", error.message);
        messageStore.setError("Something went wrong.");
      }
      console.error(error.config);
    },
  },
  persist: {
    storage: persistedState.sessionStorage,
    paths: ["companyName", "companyLogoUrl"],
  },
});
