import { api } from "@/shared/lib/api-client";
import type { GetHotTopicsResponse } from "./dto";

export const getHotTopics = async () => {
  const response = await api.get<GetHotTopicsResponse>("/topics/hot");
  return response.data;
};
