// src/shared/fixtures/problems.ts

import type { Problem } from "@/shared/api/contracts";

export function problem(
    status: number,
    title: string,
    detail?: string,
    code?: string,
    instance?: string
): Problem {
    return {
        type: "about:blank",
        title,
        status,
        detail,
        code,
        instance,
    };
}

export const problemUnauthorized = problem(401, "Unauthorized", "Missing or invalid token");
export const problemForbidden = problem(403, "Forbidden", "Not enough permissions");
export const problemNotFound = problem(404, "Not Found", "Resource not found");
export const problemNotEnoughBalance = problem(409, "Not enough balance", "amountPoints > balance", "NOT_ENOUGH_BALANCE");
export const problemValidation = problem(422, "Validation error", "phone: invalid format");
