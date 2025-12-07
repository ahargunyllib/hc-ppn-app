import {
  useGetFeedbackMetrics,
  useGetSatisfactionTrend,
} from "@/shared/repositories/feedback/query";
import { useGetHotTopics } from "@/shared/repositories/topic/query";
import { useGetUserMetrics } from "@/shared/repositories/user/query";
import {
  CircleAlertIcon,
  Clock,
  MessageSquare,
  TrendingUp,
  Users,
} from "lucide-react";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "../../shared/components/ui/alert";
import MetricCard from "./components/metric-card";
import MetricCardSkeleton from "./components/metric-card-skeleton";
import { SatisfactionLineChart } from "./components/satisfaction-line-chart";
import SatisfactionLineChartSkeleton from "./components/satisfaction-line-chart-skeleton";
import TopicsChartSkeleton from "./components/topics-chard-skeleton";
import TopicsChart from "./components/topics-chart";

export function AnalyticsDashboard() {
  const {
    data: userMetrics,
    isLoading: isLoadingUserMetrics,
    error: userMetricsError,
  } = useGetUserMetrics();

  const {
    data: feedbackMetrics,
    isLoading: isLoadingFeedbackMetrics,
    error: feedbackMetricsError,
  } = useGetFeedbackMetrics();

  const {
    data: satisfactionTrend,
    isLoading: isLoadingSatisfactionTrend,
    error: satisfactionTrendError,
  } = useGetSatisfactionTrend();

  const {
    data: hotTopics,
    isLoading: isLoadingHotTopics,
    error: hotTopicsError,
  } = useGetHotTopics();

  const isLoading =
    isLoadingUserMetrics ||
    isLoadingFeedbackMetrics ||
    isLoadingSatisfactionTrend ||
    isLoadingHotTopics;
  const error =
    userMetricsError ||
    feedbackMetricsError ||
    satisfactionTrendError ||
    hotTopicsError;
  const dataExists =
    !!userMetrics && !!feedbackMetrics && !!satisfactionTrend && !!hotTopics;

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return (
      <Alert variant="error">
        <CircleAlertIcon />
        <AlertTitle>Error loading analytics dashboard</AlertTitle>
        <AlertDescription>{error.message || "Unknown error"}</AlertDescription>
      </Alert>
    );
  }

  if (!dataExists) {
    return (
      <Alert variant="warning">
        <AlertTitle>No data available</AlertTitle>
        <AlertDescription>
          There is no analytics data to display at the moment.
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard icon={MessageSquare} title="Total Interactions" value={0} />
        <MetricCard icon={Clock} title="Avg Response Time" value={0} />
        <MetricCard
          icon={TrendingUp}
          title="Satisfaction Score"
          value={`${feedbackMetrics.payload.satisfactionScore.toFixed(1)}%`}
        />
        <MetricCard
          icon={Users}
          title="Total Users"
          value={userMetrics.payload.totalUsers.toLocaleString()}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-4">
        <div className="lg:col-span-2">
          <TopicsChart data={hotTopics.payload.topics} />
        </div>
        <div className="lg:col-span-2">
          <SatisfactionLineChart data={satisfactionTrend.payload.trend} />
        </div>
      </div>
    </div>
  );
}

function Loading() {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCardSkeleton />
        <MetricCardSkeleton />
        <MetricCardSkeleton />
        <MetricCardSkeleton />
      </div>

      <div className="grid gap-6 lg:grid-cols-4">
        <div className="lg:col-span-2">
          <TopicsChartSkeleton />
        </div>
        <div className="lg:col-span-2">
          <SatisfactionLineChartSkeleton />
        </div>
      </div>
    </div>
  );
}
