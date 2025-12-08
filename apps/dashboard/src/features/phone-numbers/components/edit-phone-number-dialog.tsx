import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogDescription,
  DialogHeader,
  DialogPanel,
  DialogPopup,
  DialogTitle,
  DialogTrigger,
} from "@/shared/components/ui/dialog";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/shared/components/ui/field";
import { Input } from "@/shared/components/ui/input";
import { toastManager } from "@/shared/components/ui/toast";
import { parseAPIError } from "@/shared/lib/api-client";
import { UpdateUserSchema } from "@/shared/repositories/user/dto";
import { useUpdateUser } from "@/shared/repositories/user/query";
import type { User } from "@/shared/types/user";
import { useForm } from "@tanstack/react-form";
import { PencilIcon } from "lucide-react";
import { useState } from "react";

type EditPhoneNumberDialogProps = {
  user: User;
};

export default function EditPhoneNumberDialog({
  user,
}: EditPhoneNumberDialogProps) {
  const [isOpen, setIsOpen] = useState(false);

  const { mutate, isPending } = useUpdateUser();

  const form = useForm({
    defaultValues: {
      phoneNumber: user.phoneNumber,
      name: user.name,
      jobTitle: user.jobTitle || "",
      gender: user.gender || "",
      dateOfBirth: user.dateOfBirth || "",
    },
    validators: {
      onSubmit: UpdateUserSchema,
    },
    onSubmit: ({ value }) => {
      mutate(
        { id: user.id, data: value },
        {
          onSuccess: () => {
            toastManager.add({
              type: "success",
              title: "User updated successfully",
            });
            setIsOpen(false);
          },
          onError: (err) => {
            toastManager.add({
              type: "error",
              title: "Failed to update user",
              description: parseAPIError(err),
            });
          },
        }
      );
    },
  });

  return (
    <Dialog onOpenChange={setIsOpen} open={isOpen}>
      <DialogTrigger render={<Button size="icon-xs" variant="ghost" />}>
        <PencilIcon />
      </DialogTrigger>
      <DialogPopup>
        <DialogHeader>
          <DialogTitle>Edit Phone Number</DialogTitle>
          <DialogDescription>
            Update the phone number and user details.
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          <form
            className="flex flex-col gap-4"
            id="edit-phone-number-form"
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <form.Field name="phoneNumber">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Phone Number</FieldLabel>
                  <FieldDescription>
                    Include country code, e.g., +62.
                  </FieldDescription>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="tel"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan nomor telepon"
                    type="tel"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="name">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Name</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="name"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Enter name"
                    type="text"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="jobTitle">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Job Title</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="organization-title"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Enter job title"
                    type="text"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="gender">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Gender</FieldLabel>
                  <select
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:font-medium file:text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    value={field.state.value}
                  >
                    <option value="">Select gender</option>
                    <option value="male">Laki-laki</option>
                    <option value="female">Perempuan</option>
                  </select>
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="dateOfBirth">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Date of Birth</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    type="date"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Subscribe>
              {() => (
                <Button
                  disabled={isPending}
                  form="edit-phone-number-form"
                  type="submit"
                >
                  {isPending ? "Loading..." : "Update"}
                </Button>
              )}
            </form.Subscribe>
          </form>
        </DialogPanel>
      </DialogPopup>
    </Dialog>
  );
}
