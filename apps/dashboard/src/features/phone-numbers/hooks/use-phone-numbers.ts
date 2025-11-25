import { mockDataStore } from "@/shared/lib/mock-data";
import type { PhoneNumber } from "@/shared/types/dashboard";
import { useQuery } from "@tanstack/react-query";

const QUERY_KEY = ["phoneNumbers"];

export function usePhoneNumbers() {
  return useQuery<PhoneNumber[]>({
    queryKey: QUERY_KEY,
    queryFn: () => mockDataStore.getPhoneNumbers(),
  });
}
