import type { User } from "./user";

export type Feedback = {
  id: string;
  rating: number;
  comment?: string;
  createdAt: string;

  user: User;
};
