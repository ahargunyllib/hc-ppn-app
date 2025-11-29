import { useInfiniteQuery } from "@tanstack/react-query";
import { getFeedbacks } from "./action";
import type { GetFeedbacksQuery } from "./dto";

export const useGetFeedbacks = (query?: GetFeedbacksQuery) =>
  useInfiniteQuery({
    queryKey: ["feedbacks", query],
    queryFn: ({ pageParam = 1 }) =>
      getFeedbacks({ ...query, page: pageParam as number }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) =>
      lastPage.payload.meta.pagination.page <
      lastPage.payload.meta.pagination.total_page
        ? lastPage.payload.meta.pagination.page + 1
        : undefined,
  });
