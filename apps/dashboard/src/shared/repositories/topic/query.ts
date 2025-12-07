import { useQuery } from "@tanstack/react-query";
import { getHotTopics } from "./action";

export const useGetHotTopics = () =>
  useQuery({
    queryKey: ["hotTopics"],
    queryFn: () => getHotTopics(),
  });
