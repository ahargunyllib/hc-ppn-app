import DataPagination from "@/shared/components/data-pagination";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/shared/components/ui/alert";
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
import { CircleAlertIcon } from "lucide-react";
import { useState } from "react";
import CreatePhoneNumberDialog from "./components/create-phone-number-dialog";
import ImportCSVDialog from "./components/import-csv-dialog";
import { PhoneNumbersTable } from "./components/phone-numbers-table";

export function PhoneNumbersManagement() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);

  const { data, isLoading, error } = useGetUsers({ page, limit });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Phone Number Management</CardTitle>
          <CardDescription>Manage user phone numbers.</CardDescription>
          <CardAction>
            <div className="flex gap-2">
              <Skeleton className="h-8 w-32 rounded-md" />
              <Skeleton className="h-8 w-32 rounded-md" />
            </div>
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
          <div className="flex gap-2">
            <ImportCSVDialog />
            <CreatePhoneNumberDialog />
          </div>
        </CardAction>
      </CardHeader>
      <CardContent>
        <PhoneNumbersTable data={data?.payload.users || []} />
        <DataPagination
          currentLimit={limit}
          currentPage={page}
          setLimit={setLimit}
          setPage={setPage}
          totalData={data?.payload.meta.pagination.total_data || 0}
          totalPage={data?.payload.meta.pagination.total_page || 1}
        />
      </CardContent>
    </Card>
  );
}
