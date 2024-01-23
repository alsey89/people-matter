import axios from "axios";

export const useJobStore = defineStore("job-store", {
  state: () => ({
    allJobs: [],
  }),
  getters: {
    //* all data
    getAllJobs: (state) => state.allJobs,
  },
  actions: {
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
          `http://localhost:3001/api/v1/job/company/${companyId}`,
          {
            companyId: jobFormData.companyId,
            title: jobFormData.title,
            departmentId: jobFormData.departmentId,
            locationId: jobFormData.locationId,
            managerId: jobFormData.managerId,
            MinSalary: jobFormData.minSalary,
            MaxSalary: jobFormData.maxSalary,
            description: jobFormData.description,
            duties: jobFormData.duties,
            qualifications: jobFormData.qualifications,
            experience: jobFormData.experience,
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
          `http://localhost:3001/api/v1/job/company/${companyId}/${jobId}`,
          {
            companyId: jobFormData.companyId,
            title: jobFormData.title,
            departmentId: jobFormData.departmentId,
            locationId: jobFormData.locationId,
            managerId: jobFormData.managerId,
            MinSalary: jobFormData.minSalary,
            MaxSalary: jobFormData.maxSalary,
            description: jobFormData.description,
            duties: jobFormData.duties,
            qualifications: jobFormData.qualifications,
            experience: jobFormData.experience,
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
          `http://localhost:3001/api/v1/job/company/${companyId}/${jobId}`,
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
      if (response.status >= 200 && response.status < 300) {
        this.allJobs = response.data.data;
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
