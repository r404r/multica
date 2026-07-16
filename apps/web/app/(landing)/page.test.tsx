import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

vi.mock(
  "@/features/landing/components/redirect-if-authenticated",
  () => ({
    RedirectIfAuthenticated: () => null,
  }),
);

import LandingPage from "./page";

describe("LandingPage", () => {
  it("renders only the login action for unauthenticated visitors", () => {
    render(<LandingPage />);

    const login = screen.getByRole("link", { name: "Login" });
    expect(login).toHaveAttribute("href", "/login");
    expect(screen.getAllByRole("link")).toHaveLength(1);
    expect(screen.queryByRole("button")).not.toBeInTheDocument();
  });
});
