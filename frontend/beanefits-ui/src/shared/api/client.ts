import { API_BASE_URL, API_DEBUG, API_TIMEOUT_MS, USE_FIXTURES } from "@/shared/config/env";

import type { Problem } from "@/shared/api/contracts";
import { mockRequestJson } from "@/shared/api/mock";
import { ApiError, NetworkError, TimeoutError } from "@/shared/api/errors";

export async function requestJson<T>(
    path: string,
    init: RequestInit & { token?: string } = {}
): Promise<T> {
    if (USE_FIXTURES) return mockRequestJson<T>(path, init);

    const url = `${API_BASE_URL}${path}`;
    const headers = new Headers(init.headers);

    headers.set("Accept", "application/json");
    if (!(init.body instanceof FormData)) headers.set("Content-Type", "application/json");
    if (init.token) headers.set("Authorization", `Bearer ${init.token}`);

    const controller = new AbortController();
    const timer = window.setTimeout(() => controller.abort(), API_TIMEOUT_MS);

    try {
        if (API_DEBUG) {
            console.log("[API] ->", init.method ?? "GET", url, {
                token: Boolean(init.token),
                body: init.body,
            });
        }
        const res = await fetch(url, { ...init, headers, signal: controller.signal });

        if (res.status === 204) return undefined as T;
        if (API_DEBUG) {
            console.log("[API] <-", res.status, url, res.headers.get("content-type"));
        }

        const ct = res.headers.get("content-type") ?? "";
        const isProblem = ct.includes("application/problem+json");

        if (!res.ok) {
            const problem: Problem = isProblem
                ? await res.json()
                : {
                    type: "about:blank",
                    title: "Request failed",
                    status: res.status,
                    detail: await res.text(),
                };

            if (problem.status === 401 && !path.startsWith("/auth/")) {
                window.dispatchEvent(new Event("auth:logout"));
            }

            throw new ApiError(problem);
        }

        return (await res.json()) as T;
    } catch (e: any) {
        if (e?.name === "AbortError") throw new TimeoutError();
        if (e instanceof ApiError) throw e;
        throw new NetworkError();
    } finally {
        window.clearTimeout(timer);
    }
}
