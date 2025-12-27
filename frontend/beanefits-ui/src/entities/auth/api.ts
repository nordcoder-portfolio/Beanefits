import type { AuthResponse } from "@/shared/api/contracts";
import { requestJson } from "@/shared/api/client";

export type LoginInput = { phone: string; password: string };
export type RegisterInput = { phone: string; password: string };

export function login(input: LoginInput): Promise<AuthResponse> {
    return requestJson<AuthResponse>("/auth/login", { method: "POST", body: JSON.stringify(input) });
}

export function register(input: RegisterInput): Promise<AuthResponse> {
    return requestJson<AuthResponse>("/auth/register", { method: "POST", body: JSON.stringify(input) });
}
