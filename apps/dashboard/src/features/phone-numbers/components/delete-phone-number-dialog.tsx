import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogPanel,
  DialogPopup,
  DialogTitle,
  DialogTrigger,
} from "@/shared/components/ui/dialog";
import { toastManager } from "@/shared/components/ui/toast";
import { parseAPIError } from "@/shared/lib/api-client";
import { useDeleteUser } from "@/shared/repositories/user/query";
import type { User } from "@/shared/types/user";
import { Trash2 } from "lucide-react";
import { useState } from "react";

type DeletePhoneNumberDialogProps = {
  user: User;
};

export default function DeletePhoneNumberDialog({
  user,
}: DeletePhoneNumberDialogProps) {
  const [open, onOpenChange] = useState(false);
  const { mutate, isPending } = useDeleteUser();

  const handleDelete = () => {
    mutate(user.id, {
      onSuccess: () => {
        toastManager.add({
          type: "success",
          title: "Phone number deleted successfully",
        });
        onOpenChange(false);
      },
      onError: (err) => {
        toastManager.add({
          type: "error",
          title: "Failed to delete phone number",
          description: parseAPIError(err),
        });
      },
    });
  };

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      <DialogTrigger render={<Button size="icon-xs" variant="ghost" />}>
        <Trash2 className="text-destructive" />
      </DialogTrigger>
      <DialogPopup>
        <DialogHeader>
          <DialogTitle>Delete Phone Number</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this phone number? This action
            cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          <div className="flex flex-col gap-2 rounded-md bg-muted p-4">
            <div className="flex justify-between">
              <span className="font-medium">Phone Number:</span>
              <span>{user.phoneNumber}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">Label:</span>
              <span>{user.label}</span>
            </div>
            {user.assignedTo && (
              <div className="flex justify-between">
                <span className="font-medium">Assigned To:</span>
                <span>{user.assignedTo}</span>
              </div>
            )}
          </div>
        </DialogPanel>
        <DialogFooter>
          <Button
            disabled={isPending}
            onClick={() => onOpenChange(false)}
            variant="ghost"
          >
            Cancel
          </Button>
          <Button
            disabled={isPending}
            onClick={handleDelete}
            variant="destructive"
          >
            {isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogPopup>
    </Dialog>
  );
}
