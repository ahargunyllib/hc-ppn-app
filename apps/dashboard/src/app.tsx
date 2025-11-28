import { AnalyticsDashboard } from "@/features/analytics/analytics-dashboard";
import { FeedbackDashboard } from "@/features/feedback/feedback-dashboard";
import { PhoneNumbersManagement } from "@/features/phone-numbers/phone-numbers-management";
import Layout from "@/shared/components/layout";
import ProtectedContent from "@/shared/components/protected-content";
import { Button } from "@/shared/components/ui/button";
import {
  Tabs,
  TabsList,
  TabsPanel,
  TabsTrigger,
} from "@/shared/components/ui/tabs";
import { useAuth } from "@/shared/hooks/use-auth";
import { BarChart3, MessageCircle, Phone } from "lucide-react";

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
    <div className="container mx-auto flex flex-col gap-6 px-4 py-8 md:py-12">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-bold text-3xl tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back, {user?.name || "User"}!
          </p>
        </div>
        <Button onClick={logout} variant="destructive">
          Logout
        </Button>
      </div>

      <Tabs defaultValue="analytics">
        <TabsList>
          <TabsTrigger value="analytics">
            <BarChart3 />
            Analytics
          </TabsTrigger>
          <TabsTrigger value="phone-numbers">
            <Phone />
            Phone Numbers
          </TabsTrigger>
          <TabsTrigger value="feedback">
            <MessageCircle />
            Feedback
          </TabsTrigger>
        </TabsList>

        <TabsPanel value="analytics">
          <AnalyticsDashboard />
        </TabsPanel>
        <TabsPanel value="phone-numbers">
          <PhoneNumbersManagement />
        </TabsPanel>
        <TabsPanel value="feedback">
          <FeedbackDashboard />
        </TabsPanel>
      </Tabs>
    </div>
  );
}
