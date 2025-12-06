import { useQuery } from "@tanstack/react-query";
import { getFeedbacks, getFeedbackMetrics, getSatisfactionTrend } from "./action";
import type { GetFeedbacksQuery, GetSatisfactionTrendQuery } from "./dto";

export const useGetFeedbacks = (query?: GetFeedbacksQuery) =>
  useQuery({
    queryKey: ["feedbacks", query],
    queryFn: () => getFeedbacks({ ...query }),
  });

export const useGetFeedbackMetrics = () =>
  useQuery({
    queryKey: ["feedbackMetrics"],
    queryFn: () => getFeedbackMetrics(),
  });

export const useGetSatisfactionTrend = (query?: GetSatisfactionTrendQuery) =>
  useQuery({
    queryKey: ["satisfactionTrend", query],
    queryFn: () => getSatisfactionTrend({ ...query }),
  });
