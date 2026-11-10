import type { Metadata } from "next";
import Link from "next/link";
import { pageMetadata } from "@/features/landing/utils/page-metadata";
import { RedirectIfAuthenticated } from "@/features/landing/components/redirect-if-authenticated";

export function generateMetadata(): Promise<Metadata> {
  return pageMetadata("home", "/");
}

export default function LandingPage() {
  return (
    <>
      <RedirectIfAuthenticated />
      <main className="flex min-h-screen items-center justify-center bg-background">
        <Link
          href="/login"
          className="inline-flex h-9 items-center justify-center rounded-lg bg-primary px-3 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/80 focus-visible:outline-none focus-visible:ring-3 focus-visible:ring-ring/50"
        >
          Login
        </Link>
      </main>
    </>
  );
}
