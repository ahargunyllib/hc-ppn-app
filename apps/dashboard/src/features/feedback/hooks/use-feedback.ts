import { apiClient } from "@/shared/lib/api-client";
import type { Feedback } from "@/shared/types/dashboard";
import { useInfiniteQuery } from "@tanstack/react-query";

const FEEDBACK_KEY = ["feedback"];

interface UseFeedbackOptions {
  userId?: string;
  minRating?: number;
  maxRating?: number;
  limit?: number;
}

export function useFeedback(options: UseFeedbackOptions = {}) {
  const { userId, minRating, maxRating, limit = 10 } = options;

  return useInfiniteQuery({
    queryKey: [...FEEDBACK_KEY, { userId, minRating, maxRating, limit }],
    queryFn: async ({ pageParam = 1 }) => {
      const response = await apiClient.getFeedbacks({
        page: pageParam,
        limit,
        userId,
        minRating,
        maxRating,
      });

      return {
        feedbacks: response.feedbacks,
        pagination: response.meta.pagination,
      };
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      const { page, total_page } = lastPage.pagination;
      return page < total_page ? page + 1 : undefined;
    },
    select: (data) => {
      const allFeedbacks = data.pages.flatMap((page) => page.feedbacks);

      const feedbacks: Feedback[] = allFeedbacks.map((feedback) => ({
        id: feedback.id,
        userId: feedback.userId,
        rating: feedback.rating,
        comment: feedback.comment,
        createdAt: feedback.createdAt,
      }));

      return {
        feedbacks,
        totalData: data.pages[0]?.pagination.total_data ?? 0,
        hasNextPage: data.pages[data.pages.length - 1]?.pagination.page < data.pages[data.pages.length - 1]?.pagination.total_page,
      };
    },
  });
}
