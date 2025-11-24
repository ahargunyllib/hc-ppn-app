import { ToastProvider } from "@/shared/components/ui/toast";
import AuthProvider from "@/shared/hooks/use-auth";

type Props = {
  children: React.ReactNode;
};

export default function Layout({ children }: Props) {
  return (
    <div className="isolate">
      <AuthProvider>
        <ToastProvider>{children}</ToastProvider>
      </AuthProvider>
    </div>
  );
}
