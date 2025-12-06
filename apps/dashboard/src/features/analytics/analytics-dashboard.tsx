import { mockDataStore } from "@/shared/lib/mock-data";
import { useGetUserMetrics } from "@/shared/repositories/user/query";
import { useQuery } from "@tanstack/react-query";
import { Clock, MessageSquare, TrendingUp, Users } from "lucide-react";
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
    isLoading: isLoadingMetrics,
    error: metricsError,
  } = useGetUserMetrics();

  const isLoading = isLoadingAnalytics || isLoadingMetrics;
  const error = analyticsError || metricsError;

  if (isLoading) {
    return <Loading />;
  }

  if (error) {
    return (
      <div className="flex h-64 items-center justify-center">
        <div className="text-center">
          <p className="text-destructive">Error loading analytics</p>
          <p className="mt-2 text-muted-foreground text-sm">
            {error instanceof Error ? error.message : "Unknown error"}
          </p>
        </div>
      </div>
    );
  }

  if (!analytics) {
    return null;
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
          value={`${metrics.satisfactionScore}%`}
        />
        <MetricCard
          icon={Users}
          title="Total Users"
          value={userMetrics?.data.totalUsers.toLocaleString() ?? "0"}
        />
      </div>

      <div className="grid gap-6 lg:grid-cols-4">
        <div className="lg:col-span-2">
          <TopicsChart data={analytics.topTopics} />
        </div>
        <div className="lg:col-span-2">
          <SatisfactionLineChart data={analytics.satisfactionTrend} />
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
