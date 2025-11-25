import { Card, CardContent } from "@/shared/components/ui/card";
import { Skeleton } from "@/shared/components/ui/skeleton";

export default function MetricCardSkeleton() {
  return (
    <Card>
      <CardContent>
        <div className="flex items-start justify-between">
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-1/2 rounded-md" />
            <div className="space-y-1">
              <Skeleton className="h-10 w-3/4 rounded-md" />
              <Skeleton className="h-3 w-1/2 rounded-md" />
            </div>
          </div>
          <div className="rounded-lg bg-muted p-2">
            <Skeleton className="size-6 text-primary" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
