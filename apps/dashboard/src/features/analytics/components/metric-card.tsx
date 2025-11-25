import { Card, CardContent } from "@/shared/components/ui/card";
import type { LucideIcon } from "lucide-react";

type MetricCardProps = {
  title: string;
  value: string | number;
  icon: LucideIcon;
};

export default function MetricCard({
  title,
  value,
  icon: Icon,
}: MetricCardProps) {
  return (
    <Card>
      <CardContent>
        <div className="flex items-start justify-between">
          <div className="flex-1 space-y-2">
            <p className="font-medium text-muted-foreground text-sm">{title}</p>
            <div className="space-y-1">
              <h3 className="font-bold text-3xl tracking-tight">{value}</h3>
            </div>
          </div>
          <div className="rounded-lg bg-muted p-2">
            <Icon className="size-6 text-primary" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
