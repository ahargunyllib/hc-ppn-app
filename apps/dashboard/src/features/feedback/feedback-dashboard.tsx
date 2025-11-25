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
import { FeedbackTable } from "./components/feedback-table";
import { useFeedback } from "./hooks/use-feedback";

export function FeedbackDashboard() {
  const {
    data: feedbacks,
    isLoading: isFeedbacksLoading,
    error: feedbacksError,
  } = useFeedback();

  if (isFeedbacksLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-1/8 rounded-md" />
          <Skeleton className="h-4 w-3/4 rounded-md" />
        </CardHeader>
        <CardContent>
          <FeedbackTable data={feedbacks || []} />
        </CardContent>
      </Card>
    );
  }

  if (feedbacksError) {
    return (
      <Alert variant="error">
        <AlertTitle>Error loading feedback</AlertTitle>
        <AlertDescription>
          {feedbacksError.message || "Unknown error"}
        </AlertDescription>
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
      <CardContent>
        <FeedbackTable data={feedbacks || []} />
      </CardContent>
    </Card>
  );
}
