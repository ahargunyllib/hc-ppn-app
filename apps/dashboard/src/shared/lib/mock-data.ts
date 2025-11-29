import type {
  AnalyticsDashboard,
  SatisfactionTrendData,
  TopicData,
} from "../types/dashboard";

// Helper function to generate random number in range
function randomInRange(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

// Helper function to format date as YYYY-MM-DD
function formatDate(date: Date): string {
  return date.toISOString().split("T")[0];
}

export function generateAnalyticsDashboard(): AnalyticsDashboard {
  const endDate = new Date();
  const startDate = new Date();
  startDate.setDate(endDate.getDate() - 30);

  return {
    metrics: {
      totalInteractions: randomInRange(5000, 15_000),
      avgResponseTime: randomInRange(30, 180),
      satisfactionScore: randomInRange(70, 95),
      activeUsers: randomInRange(500, 2000),
    },
    topTopics: generateTopTopics(),
    satisfactionTrend: generateSatisfactionTrend(30),
  };
}

function generateTopTopics(): TopicData[] {
  const topics = [
    "Account Management",
    "Billing Issues",
    "Technical Support",
    "Product Information",
    "Order Status",
    "Refund Requests",
    "Feature Requests",
    "Bug Reports",
    "General Inquiry",
    "Feedback",
  ];

  const data: TopicData[] = [];
  let remaining = 100;

  for (let i = 0; i < 5; i++) {
    const isLast = i === 4;
    const percentage = isLast ? remaining : randomInRange(10, remaining - 10);
    data.push({
      topic: topics[i],
      count: Math.floor((percentage / 100) * randomInRange(500, 1500)),
      percentage,
    });
    remaining -= percentage;
  }

  return data.sort((a, b) => b.count - a.count);
}

function generateSatisfactionTrend(days: number): SatisfactionTrendData[] {
  const data: SatisfactionTrendData[] = [];
  const baseSatisfaction = 4.2; // base average satisfaction (out of 5)

  for (let i = 0; i < days; i++) {
    const date = new Date();
    date.setDate(date.getDate() - (days - 1 - i));

    // Add random variation between -0.5 and +0.5
    const variation = (Math.random() - 0.5) * 1;
    const avgSatisfaction = Math.max(
      1,
      Math.min(5, baseSatisfaction + variation)
    );

    data.push({
      date: formatDate(date),
      avgSatisfaction: Number.parseFloat(avgSatisfaction.toFixed(2)),
    });
  }

  return data;
}

export const mockDataStore = {
  getAnalytics: (): Promise<AnalyticsDashboard> =>
    new Promise((resolve) => {
      setTimeout(() => {
        resolve(generateAnalyticsDashboard());
      }, 500);
    }),
};
