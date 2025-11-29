import type { APIResponse } from "@/shared/types/api";
import type { PaginationResponse } from "@/shared/types/pagination";
import type { User } from "@/shared/types/user";

export type GetUsersQuery = {
  page?: number;
  limit?: number;
  assignedTo?: string;
  search?: string;
};

export type GetUsersResponse = APIResponse<{
  users: User[];
  meta: {
    pagination: PaginationResponse;
  };
}>;
