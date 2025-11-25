import { mockDataStore } from "@/shared/lib/mock-data";
import type { Feedback } from "@/shared/types/dashboard";
import { useQuery } from "@tanstack/react-query";

const FEEDBACK_KEY = ["feedback"];

export function useFeedback() {
  return useQuery<Feedback[]>({
    queryKey: FEEDBACK_KEY,
    queryFn: () => mockDataStore.getFeedback(),
  });
}
