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

export type PhoneNumber = {
  id: string;
  phoneNumber: string;
  label: string;
  assignedTo?: string;
  notes?: string;
  createdAt: Date;
  updatedAt: Date;
};

export type Feedback = {
  id: string;
  sessionId: string;
  phoneNumber: string;
  rating: number; // 1-5
  comment?: string;
  createdAt: string;
};
