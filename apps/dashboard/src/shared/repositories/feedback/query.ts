import { useQuery } from "@tanstack/react-query";
import { getFeedbacks, getFeedbackMetrics } from "./action";
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
