import type { APIResponse } from "@/shared/types/api";
import type { PaginationResponse } from "@/shared/types/pagination";
import type { User } from "@/shared/types/user";
import z from "zod";

export type GetUsersQuery = {
  page?: number;
  limit?: number;
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
  name: z
    .string()
    .min(1, {
      message: "Name is required",
    })
    .max(255, {
      message: "Name must be at most 255 characters",
    }),
  jobTitle: z.union([
    z
      .string()
      .min(1, {
        message: "Job Title must be at least 1 character",
      })
      .max(255, {
        message: "Job Title must be at most 255 characters",
      }),
    z.literal(""),
  ]),
  gender: z.union([
    z.string().refine(
      (val) => {
        if (!val) {
          return true;
        }

        const allowedGenders = ["male", "female"];
        return allowedGenders.includes(val);
      },
      {
        message: "Gender must be either male or female",
      }
    ),
    z.literal(""),
  ]),
  dateOfBirth: z.union([
    z.string().refine(
      (val) => {
        if (!val) {
          return true;
        }
        const date = new Date(val);
        return !Number.isNaN(date.getTime());
      },
      {
        message: "Invalid date format",
      }
    ),
    z.literal(""),
  ]),
});

export type CreateUserRequest = z.infer<typeof CreateUserSchema>;

export type CreateUserResponse = APIResponse<{
  id: string;
}>;

export const UpdateUserSchema = z.object({
  phoneNumber: z.union([
    z
      .e164({
        error: "Invalid phone number format",
      })
      .min(10, {
        message: "Phone number must be at least 10 digits",
      })
      .max(20, {
        message: "Phone number must be at most 20 digits",
      }),
    z.literal(""),
  ]),
  name: z.union([
    z
      .string()
      .min(1, {
        message: "Name must be at least 1 character",
      })
      .max(255, {
        message: "Name must be at most 255 characters",
      }),
    z.literal(""),
  ]),
  jobTitle: z.union([
    z
      .string()
      .min(1, {
        message: "Job Title must be at least 1 character",
      })
      .max(255, {
        message: "Job Title must be at most 255 characters",
      }),
    z.literal(""),
  ]),
  gender: z.union([
    z.string().refine(
      (val) => {
        if (!val) {
          return true;
        }

        const allowedGenders = ["male", "female"];
        return allowedGenders.includes(val);
      },
      {
        message: "Gender must be either male or female",
      }
    ),
    z.literal(""),
  ]),
  dateOfBirth: z.union([
    z.string().refine(
      (val) => {
        if (!val) {
          return true;
        }
        const date = new Date(val);
        return !Number.isNaN(date.getTime());
      },
      {
        message: "Invalid date format",
      }
    ),
    z.literal(""),
  ]),
});

export type UpdateUserRequest = z.infer<typeof UpdateUserSchema>;

export type GetUserMetricsResponse = APIResponse<{
  totalUsers: number;
}>;

export type ImportUsersFromCSVResponse = APIResponse<{
  total: number;
  success: number;
  failed: number;
  errors: Array<{ row: number; error: string }>;
}>;
