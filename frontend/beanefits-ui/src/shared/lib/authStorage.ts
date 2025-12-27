import type { Account, User } from "@/shared/api/contracts";

export type AuthSnapshot = {
    accessToken: string;
    user: User;
    account: Account;
};

const KEY = "beanpoints.auth";

export function loadAuth(): AuthSnapshot | null {
    const raw = localStorage.getItem(KEY);
    if (!raw) return null;
    try {
        return JSON.parse(raw) as AuthSnapshot;
    } catch {
        return null;
    }
}

export function saveAuth(s: AuthSnapshot | null) {
    if (!s) localStorage.removeItem(KEY);
    else localStorage.setItem(KEY, JSON.stringify(s));
}
