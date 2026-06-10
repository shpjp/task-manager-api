"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { FullPageSpinner } from "@/components/ui";
import { useAuth } from "@/lib/auth-context";

export default function Home() {
  const router = useRouter();
  const { user, initializing } = useAuth();

  useEffect(() => {
    if (initializing) return;
    router.replace(user ? "/tasks" : "/login");
  }, [initializing, user, router]);

  return <FullPageSpinner />;
}
