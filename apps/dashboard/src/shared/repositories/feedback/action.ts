import { api } from "@/shared/lib/api-client";
import type { GetFeedbacksQuery, GetFeedbacksResponse, GetFeedbackMetricsResponse, GetSatisfactionTrendResponse } from "./dto";

export const getFeedbacks = async (query?: GetFeedbacksQuery) => {
  const response = await api.get<GetFeedbacksResponse>("/feedbacks", {
    params: query,
  });
  return response.data;
};

export const getFeedbackMetrics = async () => {
  const response = await api.get<GetFeedbackMetricsResponse>("/feedbacks/metrics");
  return response.data;
};

export const getSatisfactionTrend = async () => {
  const response = await api.get<GetSatisfactionTrendResponse>("/feedbacks/satisfaction-trend");
  return response.data;
};
