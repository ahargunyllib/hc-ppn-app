import DataPagination from "@/shared/components/data-pagination";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/shared/components/ui/alert";
import { Button } from "@/shared/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import { Skeleton } from "@/shared/components/ui/skeleton";
import { useGetUsers } from "@/shared/repositories/user/query";
import { CircleAlertIcon, Plus } from "lucide-react";
import { useState } from "react";
import { PhoneNumbersTable } from "./components/phone-numbers-table";

export function PhoneNumbersManagement() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);

  const {
    data,
    isLoading,
    error,
    // fetchNextPage,
    // hasNextPage,
    // isFetchingNextPage,
  } = useGetUsers({ page, limit });

  const handleAdd = () => {
    // TODO
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Phone Number Management</CardTitle>
          <CardDescription>Manage user phone numbers.</CardDescription>
          <CardAction>
            <Skeleton className="h-8 w-32 rounded-md" />
          </CardAction>
        </CardHeader>
        <CardContent>
          <PhoneNumbersTable data={[]} />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Alert variant="error">
        <CircleAlertIcon />
        <AlertTitle>Error loading phone numbers</AlertTitle>
        <AlertDescription>{error.message || "Unknown error"}</AlertDescription>
      </Alert>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Phone Number Management</CardTitle>
        <CardDescription>Manage user phone numbers.</CardDescription>
        <CardAction>
          <Button onClick={handleAdd} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Phone Number
          </Button>
        </CardAction>
      </CardHeader>
      <CardContent>
        <PhoneNumbersTable
          data={data?.pages.flatMap((p) => p.payload.users) || []}
        />
        <DataPagination
          currentLimit={limit}
          currentPage={page}
          setLimit={setLimit}
          setPage={setPage}
          totalData={data?.pages[0].payload.meta.pagination.total_data || 0}
          totalPage={data?.pages[0].payload.meta.pagination.total_page || 1}
        />
      </CardContent>
    </Card>
  );
}
