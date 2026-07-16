import type { Metadata } from "next";
import Link from "next/link";
import { buttonVariants } from "@multica/ui/components/ui/button";
import { RedirectIfAuthenticated } from "@/features/landing/components/redirect-if-authenticated";

export const metadata: Metadata = {
  title: {
    absolute: "Multica — Project Management for Human + Agent Teams",
  },
  description:
    "Open-source platform that turns coding agents into real teammates. Assign tasks, track progress, compound skills.",
  openGraph: {
    title: "Multica — Project Management for Human + Agent Teams",
    description:
      "Manage your human + agent workforce in one place.",
    url: "/",
  },
  alternates: {
    canonical: "/",
  },
};

export default function LandingPage() {
  return (
    <>
      <RedirectIfAuthenticated />
      <main className="flex min-h-screen items-center justify-center bg-background">
        <Link href="/login" className={buttonVariants({ size: "lg" })}>
          Login
        </Link>
      </main>
    </>
  );
}
