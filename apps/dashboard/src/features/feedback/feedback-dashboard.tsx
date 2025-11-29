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
import { Button } from "@/shared/components/ui/button";
import { FeedbackTable } from "./components/feedback-table";
import { useFeedback } from "./hooks/use-feedback";

export function FeedbackDashboard() {
  const {
    data,
    isLoading,
    error,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useFeedback({ limit: 10 });

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
        <AlertDescription>
          {error.message || "Unknown error"}
        </AlertDescription>
      </Alert>
    );
  }

  const feedbacks = data?.feedbacks || [];
  const totalData = data?.totalData || 0;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Feedback Dashboard</CardTitle>
        <CardDescription>
          Overview of user feedback and satisfaction metrics. Total: {totalData} feedbacks
        </CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <FeedbackTable data={feedbacks} />

        {hasNextPage && (
          <div className="flex justify-center">
            <Button
              onClick={() => fetchNextPage()}
              disabled={isFetchingNextPage}
              variant="outline"
            >
              {isFetchingNextPage ? "Loading more..." : "Load More"}
            </Button>
          </div>
        )}

        {!hasNextPage && feedbacks.length > 0 && (
          <p className="text-center text-muted-foreground text-sm">
            No more feedbacks to load
          </p>
        )}

        {feedbacks.length === 0 && (
          <p className="text-center text-muted-foreground">
            No feedbacks yet
          </p>
        )}
      </CardContent>
    </Card>
  );
}
