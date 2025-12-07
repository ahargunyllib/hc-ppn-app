import { api } from "@/shared/lib/api-client";
import type { GetHotTopicsResponse, GetTopicsCountResponse } from "./dto";

export const getHotTopics = async () => {
  const response = await api.get<GetHotTopicsResponse>("/topics/hot");
  return response.data;
};

export const getTopicsCount = async () => {
  const response = await api.get<GetTopicsCountResponse>("/topics/count");
  return response.data;
};
