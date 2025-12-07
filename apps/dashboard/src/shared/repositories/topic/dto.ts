import type { APIResponse } from "@/shared/types/api";
import type { Topic } from "@/shared/types/topic";

export type GetHotTopicsResponse = APIResponse<{
  topics: Topic[];
}>;

export type GetTopicsCountResponse = APIResponse<{
  totalTopics: number;
}>;
