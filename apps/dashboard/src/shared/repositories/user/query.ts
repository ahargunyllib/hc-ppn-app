import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createUser, deleteUser, getUsers } from "./action";
import type { CreateUserRequest, GetUsersQuery } from "./dto";

export const useGetUsers = (query?: GetUsersQuery) =>
  useQuery({
    queryKey: ["users", query],
    queryFn: () => getUsers({ ...query }),
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
