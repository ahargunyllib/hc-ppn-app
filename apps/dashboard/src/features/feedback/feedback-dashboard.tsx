import DataPagination from "@/shared/components/data-pagination";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/shared/components/ui/alert";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import { Skeleton } from "@/shared/components/ui/skeleton";
import { useGetFeedbacks } from "@/shared/repositories/feedback/query";
import { useState } from "react";
import { FeedbackTable } from "./components/feedback-table";

export function FeedbackDashboard() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);

  const { data, isLoading, error } = useGetFeedbacks({ page, limit });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-1/8 rounded-md" />
          <Skeleton className="h-4 w-3/4 rounded-md" />
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-2">
            <Skeleton className="h-12 w-full rounded-md" />
            <Skeleton className="h-12 w-full rounded-md" />
            <Skeleton className="h-12 w-full rounded-md" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Alert variant="error">
        <AlertTitle>Error loading feedback</AlertTitle>
        <AlertDescription>{error.message || "Unknown error"}</AlertDescription>
      </Alert>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Feedback Dashboard</CardTitle>
        <CardDescription>
          Overview of user feedback and satisfaction metrics.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <FeedbackTable data={data?.payload.feedbacks || []} />
        <DataPagination
          currentLimit={limit}
          currentPage={page}
          setLimit={setLimit}
          setPage={setPage}
          totalData={data?.payload.meta.pagination.total_data || 0}
          totalPage={data?.payload.meta.pagination.total_page || 1}
        />
      </CardContent>
    </Card>
  );
}
