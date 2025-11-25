import type {
  AnalyticsDashboard,
  Feedback,
  PhoneNumber,
  SatisfactionTrendData,
  TopicData,
} from "../types/dashboard";

// Helper function to generate random number in range
function randomInRange(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

// Helper function to generate random date in past N days
function randomDateInPastDays(days: number): Date {
  const date = new Date();
  date.setDate(date.getDate() - randomInRange(0, days));
  date.setHours(randomInRange(0, 23), randomInRange(0, 59), 0, 0);
  return date;
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

const PHONE_NUMBER_PREFIXES = [
  "+1-555",
  "+1-444",
  "+1-666",
  "+1-777",
  "+1-888",
];
const PHONE_LABELS = [
  "Customer Support",
  "Sales Team",
  "Technical Support",
  "Billing Department",
  "General Inquiry",
  "Emergency Line",
  "VIP Support",
  "Regional Office",
  "Marketing",
  "HR Department",
];

const ASSIGNED_USERS = [
  "John Doe",
  "Jane Smith",
  "Mike Johnson",
  "Sarah Williams",
  "David Brown",
  "Emily Davis",
  "Chris Wilson",
  "Amanda Taylor",
  "Not Assigned",
];

export function generatePhoneNumbers(count = 50): PhoneNumber[] {
  const phoneNumbers: PhoneNumber[] = [];

  for (let i = 0; i < count; i++) {
    const createdAt = randomDateInPastDays(365);
    const updatedAt = new Date(
      createdAt.getTime() + randomInRange(0, 30) * 24 * 60 * 60 * 1000
    );

    phoneNumbers.push({
      id: `phone-${i + 1}`,
      phoneNumber: `${PHONE_NUMBER_PREFIXES[randomInRange(0, PHONE_NUMBER_PREFIXES.length - 1)]}-${randomInRange(100, 999)}-${randomInRange(1000, 9999)}`,
      label: `${PHONE_LABELS[randomInRange(0, PHONE_LABELS.length - 1)]} ${randomInRange(1, 5)}`,
      assignedTo:
        Math.random() > 0.3
          ? ASSIGNED_USERS[randomInRange(0, ASSIGNED_USERS.length - 2)]
          : undefined,
      notes:
        Math.random() > 0.5
          ? `Notes for phone number ${i + 1}. Last updated on ${formatDate(updatedAt)}.`
          : undefined,
      createdAt,
      updatedAt,
    });
  }

  return phoneNumbers.sort(
    (a, b) => b.createdAt.getTime() - a.createdAt.getTime()
  );
}

const FEEDBACK_COMMENTS = {
  positive: [
    "Great service! Very satisfied with the quick response time.",
    "The support team was extremely helpful and resolved my issue promptly.",
    "Excellent experience. The interface is intuitive and easy to use.",
    "Fast response and professional handling of my request.",
    "Very pleased with the quality of service provided.",
  ],
  neutral: [
    "The service was okay. Room for improvement in response time.",
    "It works but could be better. Some features are missing.",
    "Average experience. Nothing special but gets the job done.",
    "Decent service, though I expected faster resolution.",
    "The system works fine but the UI could be more modern.",
  ],
  negative: [
    "Response time was too slow. Waited over an hour for support.",
    "The interface is confusing and not user-friendly.",
    "Had issues with the service. Expected better quality.",
    "Disappointed with the support. Issue still not fully resolved.",
    "System keeps crashing. Needs urgent fixing.",
  ],
};

export function generateFeedback(count = 150): Feedback[] {
  const feedback: Feedback[] = [];

  for (let i = 0; i < count; i++) {
    const rating = randomInRange(1, 5);
    const createdAt = randomDateInPastDays(90);

    let commentPool: string[];
    if (rating >= 4) {
      commentPool = FEEDBACK_COMMENTS.positive;
    } else if (rating === 3) {
      commentPool = FEEDBACK_COMMENTS.neutral;
    } else {
      commentPool = FEEDBACK_COMMENTS.negative;
    }

    feedback.push({
      id: `feedback-${i + 1}`,
      phoneNumber: {
        id: `phone-${randomInRange(1, 50)}`,
        phoneNumber: `${PHONE_NUMBER_PREFIXES[randomInRange(0, PHONE_NUMBER_PREFIXES.length - 1)]}-${randomInRange(100, 999)}-${randomInRange(1000, 9999)}`,
        label: PHONE_LABELS[randomInRange(0, PHONE_LABELS.length - 1)],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      rating,
      comment: commentPool[randomInRange(0, commentPool.length - 1)],
      createdAt,
    });
  }

  return feedback.sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());
}

let phoneNumbersStore = generatePhoneNumbers(10);
let feedbackStore = generateFeedback(10);

export const mockDataStore = {
  getAnalytics: (): Promise<AnalyticsDashboard> =>
    new Promise((resolve) => {
      setTimeout(() => {
        resolve(generateAnalyticsDashboard());
      }, 500);
    }),

  getPhoneNumbers: (): Promise<PhoneNumber[]> =>
    new Promise((resolve) => {
      setTimeout(() => {
        resolve([...phoneNumbersStore]);
      }, 400);
    }),

  getFeedback: (): Promise<Feedback[]> =>
    new Promise((resolve) => {
      setTimeout(() => {
        resolve([...feedbackStore]);
      }, 400);
    }),

  resetStores: (): void => {
    phoneNumbersStore = generatePhoneNumbers(50);
    feedbackStore = generateFeedback(150);
  },
};
