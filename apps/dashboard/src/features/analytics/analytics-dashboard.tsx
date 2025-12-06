import { mockDataStore } from "@/shared/lib/mock-data";
import { useGetFeedbackMetrics, useGetSatisfactionTrend } from "@/shared/repositories/feedback/query";
import { useGetUserMetrics } from "@/shared/repositories/user/query";
import { useQuery } from "@tanstack/react-query";
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
    data: analytics,
    isLoading: isLoadingAnalytics,
    error: analyticsError,
  } = useQuery({
    queryKey: ["analytics"],
    queryFn: () => mockDataStore.getAnalytics(),
  });

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
  } = useGetSatisfactionTrend({ days: 30 });

  const isLoading = isLoadingAnalytics || isLoadingUserMetrics || isLoadingFeedbackMetrics || isLoadingSatisfactionTrend;
  const error = analyticsError || userMetricsError || feedbackMetricsError || satisfactionTrendError;
  const dataExists = !!analytics && !!userMetrics && !!feedbackMetrics && !!satisfactionTrend;

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

  const { metrics } = analytics;

  const formatResponseTime = (seconds: number): string => {
    if (seconds < 60) {
      return `${seconds}s`;
    }
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}m ${remainingSeconds}s`;
  };

  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          icon={MessageSquare}
          title="Total Interactions"
          value={metrics.totalInteractions.toLocaleString()}
        />
        <MetricCard
          icon={Clock}
          title="Avg Response Time"
          value={formatResponseTime(metrics.avgResponseTime)}
        />
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
          <TopicsChart data={analytics.topTopics} />
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
