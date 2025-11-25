/**
 * Login Form Component
 * Handles user authentication with email and password
 */

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import { AUTH_CONSTANTS } from "@/shared/constant/auth";

import { Button } from "@/shared/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field";
import { Input } from "@/shared/components/ui/input";
import { useAuth } from "@/shared/hooks/use-auth";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { toastManager } from "../../../shared/components/ui/toast";

export default function LoginForm() {
  const auth = useAuth();

  const form = useForm({
    defaultValues: {
      email: "",
      password: "",
    },
    validators: {
      onSubmit: z.object({
        email: z.email("Email tidak valid"),
        password: z
          .string()
          .min(1, "Password wajib diisi")
          .min(
            AUTH_CONSTANTS.MIN_PASSWORD_LENGTH,
            "Password minimal 8 karakter"
          ),
      }),
    },
    onSubmit: ({ value }) => {
      toastManager.promise(auth.login(value), {
        loading: "Memproses login...",
        success: "Berhasil login!",
        error: (err) => (err instanceof Error ? err.message : "Gagal login"),
      });
    },
  });

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Selamat Datang Kembali!</CardTitle>
          <CardDescription>
            Silakan masuk ke akun Anda untuk melanjutkan.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form
            className="flex flex-col gap-4"
            id="login-form"
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
          >
            <form.Field name="email">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Email</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="username"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan email Anda"
                    type="email"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Field name="password">
              {(field) => (
                <Field
                  dirty={field.state.meta.isDirty}
                  invalid={!field.state.meta.isValid}
                  name={field.name}
                  touched={field.state.meta.isTouched}
                >
                  <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                  <Input
                    aria-invalid={
                      field.state.meta.isTouched && !field.state.meta.isValid
                    }
                    autoComplete="current-password"
                    id={field.name}
                    name={field.name}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Masukkan password Anda"
                    type="password"
                    value={field.state.value}
                  />
                  <FieldError match={!field.state.meta.isValid}>
                    {field.state.meta.errors[0]?.message}
                  </FieldError>
                </Field>
              )}
            </form.Field>

            <form.Subscribe>
              {(state) => (
                <Button
                  disabled={!state.canSubmit || state.isSubmitting}
                  form="login-form"
                  type="submit"
                >
                  {state.isSubmitting ? "Loading..." : "Login"}
                </Button>
              )}
            </form.Subscribe>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
