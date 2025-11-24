import Layout from "./shared/components/layout";
import ProtectedContent from "./shared/components/protected-content";
import { Button } from "./shared/components/ui/button";
import { useAuth } from "./shared/hooks/use-auth";

export default function App() {
  return (
    <Layout>
      <ProtectedContent>
        <Dashboard />
      </ProtectedContent>
    </Layout>
  );
}

function Dashboard() {
  const { user, logout } = useAuth();

  return (
    <div className="container mx-auto flex flex-col gap-6 py-20">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-bold text-3xl">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back, {user?.name || "User"}!
          </p>
        </div>
        <Button onClick={logout} variant="outline">
          Logout
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div className="rounded-lg border bg-card p-6">
          <h2 className="mb-2 font-semibold text-lg">Overview</h2>
          <p className="text-muted-foreground text-sm">
            Your dashboard content goes here
          </p>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <h2 className="mb-2 font-semibold text-lg">Statistics</h2>
          <p className="text-muted-foreground text-sm">
            View your analytics and metrics
          </p>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <h2 className="mb-2 font-semibold text-lg">Settings</h2>
          <p className="text-muted-foreground text-sm">
            Manage your account settings
          </p>
        </div>
      </div>
    </div>
  );
}
