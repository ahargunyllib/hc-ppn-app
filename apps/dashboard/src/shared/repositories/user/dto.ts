import type { APIResponse } from "@/shared/types/api";
import type { PaginationResponse } from "@/shared/types/pagination";
import type { User } from "@/shared/types/user";
import z from "zod";

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

export const CreateUserSchema = z.object({
  phoneNumber: z
    .e164({
      error: "Invalid phone number format",
    })
    .min(10, {
      message: "Phone number must be at least 10 digits",
    })
    .max(20, {
      message: "Phone number must be at most 20 digits",
    }),
  label: z
    .string()
    .min(1, {
      message: "Label is required",
    })
    .max(255, {
      message: "Label must be at most 255 characters",
    }),
  // https://stackoverflow.com/questions/73715295/react-hook-form-with-zod-resolver-optional-field
  assignedTo: z.union([
    z
      .string()
      .min(1, {
        message: "Assigned To must be at least 1 character",
      })
      .max(255, {
        message: "Assigned To must be at most 255 characters",
      }),
    z.literal(""),
  ]),
  notes: z.union([
    z.string().max(1000, {
      message: "Notes must be at most 1000 characters",
    }),
    z.literal(""),
  ]),
});

export type CreateUserRequest = z.infer<typeof CreateUserSchema>;

export type CreateUserResponse = APIResponse<{
  id: string;
}>;
