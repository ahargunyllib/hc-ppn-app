import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart";
import type { SatisfactionTrendData } from "@/shared/types/dashboard";
import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts";

type SatisfactionLineChartProps = {
  data: SatisfactionTrendData[];
};

export function SatisfactionLineChart({ data }: SatisfactionLineChartProps) {
  const chartConfig = {
    avgSatisfaction: {
      label: "Avg Satisfaction",
      color: "var(--chart-1)",
    },
  };

  // Format date for display (show only last 15 entries for readability)
  const displayData = data.slice(-15).map((item) => ({
    ...item,
    date: new Date(item.date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    }),
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle>Satisfaction Trend</CardTitle>
        <CardDescription>
          Average user satisfaction ratings over time
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <LineChart
            data={displayData}
            margin={{ top: 5, right: 5, left: 5, bottom: 5 }}
          >
            <CartesianGrid className="stroke-muted" strokeDasharray="3 3" />
            <XAxis className="text-xs" dataKey="date" />
            <YAxis
              className="text-xs"
              domain={[1, 5]}
              ticks={[1, 2, 3, 4, 5]}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend content={<ChartLegendContent />} />
            <Line
              dataKey="avgSatisfaction"
              dot={false}
              stroke="var(--color-avgSatisfaction)"
              strokeWidth={2}
              type="monotone"
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
