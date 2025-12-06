import { api } from "@/shared/lib/api-client";
import type {
  CreateUserRequest,
  GetUserMetricsResponse,
  GetUsersQuery,
  GetUsersResponse,
} from "./dto";

export const getUsers = async (query?: GetUsersQuery) => {
  const response = await api.get<GetUsersResponse>("/users", { params: query });
  return response.data;
};

export const createUser = async (req: CreateUserRequest) => {
  const response = await api.post("/users", req);
  return response.data;
};

export const deleteUser = async (id: string) => {
  await api.delete(`/users/${id}`);
};

export const getUserMetrics = async () => {
  const response = await api.get<GetUserMetricsResponse>("/users/metrics");
  return response.data;
};
