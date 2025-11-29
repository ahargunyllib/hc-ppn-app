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
import { Textarea } from "@/shared/components/ui/textarea";
import { toastManager } from "@/shared/components/ui/toast";
import { parseAPIError } from "@/shared/lib/api-client";
import { CreateUserSchema } from "@/shared/repositories/user/dto";
import { useCreateUser } from "@/shared/repositories/user/query";
import { useForm } from "@tanstack/react-form";
import { Plus } from "lucide-react";
import { useState } from "react";

export default function CreatePhoneNumberDialog() {
  const [isOpen, setIsOpen] = useState(false);

  const { mutate, isPending } = useCreateUser();

  const form = useForm({
    defaultValues: {
      phoneNumber: "",
      label: "",
      assignedTo: "",
      notes: "",
    },
    validators: {
      onSubmit: CreateUserSchema,
    },
    onSubmit: ({ value }) => {
      mutate(value, {
        onSuccess: () => {
          toastManager.add({
            type: "success",
            title: "User created successfully",
          });
          setIsOpen(false);
          form.reset();
        },
        onError: (err) => {
          toastManager.add({
            type: "error",
            title: "Failed to create user",
            description: parseAPIError(err),
          });
        },
      });
    },
  });

  return (
    <Dialog onOpenChange={setIsOpen} open={isOpen}>
      <DialogTrigger render={<Button size="sm" />}>
        <Plus className="mr-2 h-4 w-4" />
        Add Phone Number
      </DialogTrigger>
      <DialogPopup>
        <DialogHeader>
          <DialogTitle>Create Phone Number</DialogTitle>
          <DialogDescription>
            Create a new phone number for a user.
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          <form
            className="flex flex-col gap-4"
            id="create-phone-number-form"
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
                  <FieldLabel htmlFor={field.name}>
                    Phone Number <span className="text-red-500">*</span>
                  </FieldLabel>
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

            <form.Field name="label">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>
                    Label
                    <span className="text-red-500">*</span>
                  </FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="off"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan label"
                    type="text"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="assignedTo">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Assigned To</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="off"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan user ID yang ditugaskan"
                    type="text"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="notes">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Notes</FieldLabel>
                  <Textarea
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="off"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan catatan"
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
                  form="create-phone-number-form"
                  type="submit"
                >
                  {isPending ? "Loading..." : "Create"}
                </Button>
              )}
            </form.Subscribe>
          </form>
        </DialogPanel>
      </DialogPopup>
    </Dialog>
  );
}
