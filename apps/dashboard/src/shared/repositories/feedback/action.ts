import { api } from "@/shared/lib/api-client";
import type { GetFeedbacksQuery, GetFeedbacksResponse } from "./dto";

export const getFeedbacks = async (query?: GetFeedbacksQuery) => {
  const response = await api.get<GetFeedbacksResponse>("/feedbacks", {
    params: query,
  });
  return response.data;
};
