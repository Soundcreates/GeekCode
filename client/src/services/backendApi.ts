import axios from "axios";

export const fetchData = axios.create({
  baseURL: "http://localhost:8080/api",
});

// Add request interceptor to include token dynamically
fetchData.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
