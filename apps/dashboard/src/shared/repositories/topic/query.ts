import { useQuery } from "@tanstack/react-query";
import { getHotTopics, getTopicsCount } from "./action";

export const useGetHotTopics = () =>
  useQuery({
    queryKey: ["hotTopics"],
    queryFn: () => getHotTopics(),
  });

export const useGetTopicsCount = () =>
  useQuery({
    queryKey: ["topicsCount"],
    queryFn: () => getTopicsCount(),
  });
