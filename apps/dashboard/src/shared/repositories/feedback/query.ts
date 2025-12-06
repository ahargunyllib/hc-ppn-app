import { useQuery } from "@tanstack/react-query";
import { getFeedbacks, getFeedbackMetrics, getSatisfactionTrend } from "./action";
import type { GetFeedbacksQuery } from "./dto";

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

export const useGetSatisfactionTrend = () =>
  useQuery({
    queryKey: ["satisfactionTrend"],
    queryFn: () => getSatisfactionTrend(),
  });
