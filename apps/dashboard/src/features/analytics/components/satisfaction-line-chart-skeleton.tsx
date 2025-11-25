import { Card, CardContent, CardHeader } from "@/shared/components/ui/card";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart";
import { Skeleton } from "@/shared/components/ui/skeleton";
import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts";

export default function SatisfactionLineChartSkeleton() {
  const chartConfig = {
    avgSatisfaction: {
      label: "Avg Satisfaction",
      color: "var(--chart-1)",
    },
  };

  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-6 w-1/4 rounded-md" />
        <Skeleton className="h-4 w-3/4 rounded-md" />
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <LineChart
            data={[{}, {}, {}, {}, {}]}
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
