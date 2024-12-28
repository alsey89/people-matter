import axios from "axios";

const baseUrl = import.meta.env.VITE_API_URL;

const instance = axios.create({
  baseURL: baseUrl,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
  withCredentials: true,
  timeout: 50000,
});

instance.interceptors.request.use(
  (config) => {
    // Do something before request is sent
    return config;
  },
  (error) => {
    // Do something before request is sent
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  (response) => {
    // Do something with response data
    return response;
  },
  (error) => {
    // if (error.response) {
    //   switch (error.response.status) {
    //     case 400:
    //       console.log("bad request");
    //       break;
    //     case 404:
    //       console.log("not found");
    //       // go to 404 page
    //       break;
    //     case 500:
    //       console.log("something went wrong");
    //       // go to 500 page
    //       break;
    //     default:
    //       console.log(error.message);
    //   }
    // }
    if (!window.navigator.onLine) {
      console.error("No internet connection. Please try again later.");
      return;
    }
    return Promise.reject(error);
  }
);

export default instance;
