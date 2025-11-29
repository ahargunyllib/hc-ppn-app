import {
  useInfiniteQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { createUser, deleteUser, getUsers } from "./action";
import type { CreateUserRequest, GetUsersQuery } from "./dto";

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

export const useCreateUser = () => {
  const qc = useQueryClient();

  return useMutation({
    mutationKey: ["createUser"],
    mutationFn: (req: CreateUserRequest) => createUser(req),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
    },
  });
};

export const useDeleteUser = () => {
  const qc = useQueryClient();

  return useMutation({
    mutationKey: ["deleteUser"],
    mutationFn: (id: string) => deleteUser(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
    },
  });
};
