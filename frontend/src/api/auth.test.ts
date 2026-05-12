import { beforeEach, describe, expect, it, vi } from "vitest";
import { getSession } from "./auth";

describe("getSession", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it("throw an error when the response returns as !ok", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 401
      })
    );

    await expect(getSession()).rejects.toThrow("session check failed: 401");
  });

  it("returns the parsed session with a success status", async () => {
    const mockSession = {
      userId: "123",
      username: "testuser",
      createdAt: "2024-01-01T00:00:00Z",
      status: "active",
      lastSeenAt: "2024-01-01T12:00:00Z",
      avatarUrl: ""
    };

    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: true,
        json: vi.fn().mockResolvedValue(mockSession)
      })
    );

    localStorage.setItem("cabin_token", "valid_token");
    const session = await getSession();
    expect(session.username).toBe("testuser");
  });
});
