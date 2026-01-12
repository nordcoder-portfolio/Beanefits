import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import type { Account, RoleCode, User } from "@/shared/api/contracts";
import { loadAuth, saveAuth } from "@/shared/lib/authStorage";
import * as authApi from "@/entities/auth/api";

type AuthState = {
    accessToken: string | null;
    user: User | null;
    account: Account | null;
};

type AuthContextValue = AuthState & {
    login: (phone: string, password: string) => Promise<void>;
    register: (phone: string, password: string) => Promise<void>;
    logout: () => void;
    hasRole: (role: RoleCode) => boolean;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const initial = loadAuth();
    const [state, setState] = useState<AuthState>(() => ({
        accessToken: initial?.accessToken ?? null,
        user: initial?.user ?? null,
        account: initial?.account ?? null,
    }));

    const commit = useCallback((s: AuthState) => {
        setState(s);
        if (s.accessToken && s.user && s.account) {
            saveAuth({ accessToken: s.accessToken, user: s.user, account: s.account });
        } else {
            saveAuth(null);
        }
    }, []);

    const login = useCallback(async (phone: string, password: string) => {
        const res = await authApi.login({ phone, password });
        commit({ accessToken: res.accessToken, user: res.user, account: res.account });
        return res.user;
    }, [commit]);

    const register = useCallback(async (phone: string, password: string) => {
        const res = await authApi.register({ phone, password });
        commit({ accessToken: res.accessToken, user: res.user, account: res.account });
        return res.user;
    }, [commit]);

    const logout = useCallback(() => commit({ accessToken: null, user: null, account: null }), [commit]);
    useEffect(() => {
        const onLogout = () => logout();
        window.addEventListener("auth:logout", onLogout);
        return () => window.removeEventListener("auth:logout", onLogout);
    }, [logout]);

    const hasRole = useCallback((role: RoleCode) => {
        return Boolean(state.user?.roles?.includes(role));
    }, [state.user]);

    const value = useMemo<AuthContextValue>(() => ({ ...state, login, register, logout, hasRole }), [state, login, register, logout, hasRole]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error("useAuth must be used within AuthProvider");
    return ctx;
}
