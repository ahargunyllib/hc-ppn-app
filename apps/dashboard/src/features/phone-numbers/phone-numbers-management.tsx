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
import { CircleAlertIcon, Plus } from "lucide-react";
import { PhoneNumbersTable } from "./components/phone-numbers-table";
import { usePhoneNumbers } from "./hooks/use-phone-numbers";

export function PhoneNumbersManagement() {
  const { data: phoneNumbers, isLoading, error } = usePhoneNumbers();

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
        <PhoneNumbersTable data={phoneNumbers || []} />
      </CardContent>
    </Card>
  );
}
