import type { Problem } from "@/shared/api/contracts";

export class ApiError extends Error {
    constructor(public problem: Problem) {
        super(problem.title);
    }
}

export class NetworkError extends Error {
    constructor(message = "Cannot reach server") {
        super(message);
    }
}

export class TimeoutError extends Error {
    constructor(message = "Server did not respond in time") {
        super(message);
    }
}

export function isApiError(e: unknown): e is ApiError {
    return e instanceof ApiError && typeof e.problem?.status === "number";
}

export function isNetworkError(e: unknown): e is NetworkError {
    return e instanceof NetworkError;
}

export function isTimeoutError(e: unknown): e is TimeoutError {
    return e instanceof TimeoutError;
}

/**
 * Единое человекочитаемое сообщение для UI.
 * Никаких “пароль неверный”, если сервер не доступен.
 */
export function humanizeError(e: unknown): string {
    if (isTimeoutError(e)) return "Server is not responding. Please try again.";
    if (isNetworkError(e)) return "Cannot connect to server. Check that backend is running and the URL is correct.";

    if (isApiError(e)) {
        const p = e.problem;

        if (p.status === 401) return p.detail || "Invalid credentials.";
        if (p.status === 403) return p.detail || "Not enough permissions.";
        if (p.status === 409) return p.detail || "Conflict.";
        if (p.status === 422) return p.detail || "Validation error.";

        if (p.status >= 500) return "Server error. Please try again.";
        return p.detail || p.title || "Request failed.";
    }

    const msg = (e as any)?.message;
    return typeof msg === "string" && msg.trim() ? msg : "Unexpected error.";
}
