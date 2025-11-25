import { Card, CardContent, CardHeader } from "@/shared/components/ui/card";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart";
import { Skeleton } from "@/shared/components/ui/skeleton";
import { Bar, BarChart, XAxis, YAxis } from "recharts";

export default function TopicsChartSkeleton() {
  const chartConfig = {
    count: {
      label: "Interactions",
      color: "var(--chart-1)",
    },
  };

  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-6 w-1/8 rounded-md" />
        <Skeleton className="h-4 w-3/4 rounded-md" />
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart data={[{}, {}, {}, {}, {}, {}]}>
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
