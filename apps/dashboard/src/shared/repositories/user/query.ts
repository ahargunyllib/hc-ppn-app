import { useInfiniteQuery } from "@tanstack/react-query";
import { getUsers } from "./action";
import type { GetUsersQuery } from "./dto";

export const useGetUsers = (query?: GetUsersQuery) =>
  useInfiniteQuery({
    queryKey: ["users", query],
    queryFn: ({ pageParam = 1 }) =>
      getUsers({ ...query, page: pageParam as number }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) =>
      lastPage.payload.meta.pagination.page <
      lastPage.payload.meta.pagination.total_page
        ? lastPage.payload.meta.pagination.page + 1
        : undefined,
  });
