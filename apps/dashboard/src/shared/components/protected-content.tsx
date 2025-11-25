/**
 * Protected Content Component
 * Wrapper that shows content only when user is authenticated
 */

import LoginForm from "@/features/auth/components/login-form";
import { Spinner } from "@/shared/components/ui/spinner";
import { useAuth } from "@/shared/hooks/use-auth";
import type { ReactNode } from "react";

type ProtectedContentProps = {
  children: ReactNode;
};

/**
 * Renders children only if user is authenticated
 * Shows login form by default, or custom fallback if provided
 */
export default function ProtectedContent({ children }: ProtectedContentProps) {
  const { isAuthenticated, isLoading } = useAuth();

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center">
          <Spinner />
          <p className="text-muted-foreground text-sm">Loading...</p>
        </div>
      </div>
    );
  }

  // Show login form if not authenticated
  if (!isAuthenticated) {
    return <LoginForm />;
  }

  // Show protected content
  return <>{children}</>;
}
