import { beforeEach, describe, expect, it } from "vitest";
import { authHeaders } from "./headers";

describe("authHeaders", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("returns only Content-Type when there is no token present", () => {
    const headers = authHeaders() as Record<string, string>;
    expect(headers["Content-Type"]).toBe("application/json");
    expect(headers["Authorization"]).toBeUndefined();
  });

  it("includes a Bearer token when a token is present", () => {
    localStorage.setItem("cabin_token", "testtoken");
    const headers = authHeaders() as Record<string, string>;
    expect(headers["Content-Type"]).toBe("application/json");
    expect(headers["Authorization"]).toBe("Bearer testtoken");
  });
});
