import { ToastProvider } from "@/shared/components/ui/toast";

type Props = {
  children: React.ReactNode;
};

export default function Layout({ children }: Props) {
  return (
    <div className="isolate">
      <ToastProvider>{children}</ToastProvider>
    </div>
  );
}
