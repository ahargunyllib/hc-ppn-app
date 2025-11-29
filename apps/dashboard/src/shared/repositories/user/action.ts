import { api } from "@/shared/lib/api-client";
import type { GetUsersQuery, GetUsersResponse } from "./dto";

export const getUsers = async (query?: GetUsersQuery) => {
  const response = await api.get<GetUsersResponse>("/users", { params: query });
  return response.data;
};
