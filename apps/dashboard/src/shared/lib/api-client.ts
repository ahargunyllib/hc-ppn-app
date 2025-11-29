import axios from "axios";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export const api = axios.create({
  baseURL: API_BASE_URL,
});

export const parseAPIError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    return error.response?.data.payload.error.message || error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return String(error);
};
