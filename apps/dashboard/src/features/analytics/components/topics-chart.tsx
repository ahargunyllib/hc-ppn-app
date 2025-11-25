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
import type { TopicData } from "@/shared/types/dashboard";
import { Bar, BarChart, XAxis, YAxis } from "recharts";

type TopicsChartProps = {
  data: TopicData[];
};

export default function TopicsChart({ data }: TopicsChartProps) {
  const chartConfig = {
    count: {
      label: "Interactions",
      color: "var(--chart-1)",
    },
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Popular Topics</CardTitle>
        <CardDescription>
          Topics with the highest user interactions over the selected period
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart data={data}>
            <XAxis className="text-xs" dataKey="topic" />
            <YAxis className="text-xs" />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend content={<ChartLegendContent />} />
            <Bar
              dataKey="count"
              fill="var(--color-count)"
              radius={[8, 8, 0, 0]}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
