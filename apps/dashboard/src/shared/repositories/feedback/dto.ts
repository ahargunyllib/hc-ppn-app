import type { APIResponse } from "@/shared/types/api";
import type { Feedback } from "@/shared/types/feedback";
import type { PaginationResponse } from "@/shared/types/pagination";

export type GetFeedbacksQuery = {
  page?: number;
  limit?: number;
  userId?: string;
  minRating?: number;
  maxRating?: number;
};

export type GetFeedbacksResponse = APIResponse<{
  feedbacks: Feedback[];
  meta: {
    pagination: PaginationResponse;
  };
}>;

export type GetFeedbackMetricsResponse = APIResponse<{
  satisfactionScore: number;
}>;
