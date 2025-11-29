export type InteractionMetrics = {
  totalInteractions: number;
  avgResponseTime: number; // in seconds
  satisfactionScore: number; // 0-100
  activeUsers: number;
};

export type TopicData = {
  topic: string;
  count: number;
  percentage: number;
};

export type SatisfactionTrendData = {
  date: string;
  avgSatisfaction: number; // 1-5
};

export type AnalyticsDashboard = {
  metrics: InteractionMetrics;
  topTopics: TopicData[];
  satisfactionTrend: SatisfactionTrendData[];
};
