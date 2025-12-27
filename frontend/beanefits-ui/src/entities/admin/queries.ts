import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/app/providers/AuthProvider";
import type { CreateRulesetRequest } from "@/shared/api/contracts";

import {
    listUsers,
    deleteUser,
    getCurrentRuleset,
    listRulesets,
    createRuleset,
} from "./api";

export function useAdminUsers(params: { limit?: number; offset?: number; q?: string }) {
    const { accessToken } = useAuth();

    const limit = params.limit ?? 20;
    const offset = params.offset ?? 0;
    const q = params.q ?? "";

    return useQuery({
        queryKey: ["admin", "users", limit, offset, q],
        queryFn: () => listUsers(accessToken!, { limit, offset, q: q || undefined }),
        enabled: Boolean(accessToken),
    });
}

export function useDeleteUser() {
    const { accessToken } = useAuth();
    const qc = useQueryClient();

    return useMutation({
        mutationFn: (userId: number) => deleteUser(accessToken!, userId),
        onSuccess: () => qc.invalidateQueries({ queryKey: ["admin", "users"] }),
    });
}

export function useCurrentRuleset() {
    const { accessToken } = useAuth();

    return useQuery({
        queryKey: ["admin", "rulesets", "current"],
        queryFn: () => getCurrentRuleset(accessToken!),
        enabled: Boolean(accessToken),
    });
}

export function useRulesets(params: { limit?: number; offset?: number }) {
    const { accessToken } = useAuth();
    const limit = params.limit ?? 20;
    const offset = params.offset ?? 0;

    return useQuery({
        queryKey: ["admin", "rulesets", limit, offset],
        queryFn: () => listRulesets(accessToken!, { limit, offset }),
        enabled: Boolean(accessToken),
    });
}

export function useCreateRuleset() {
    const { accessToken } = useAuth();
    const qc = useQueryClient();

    return useMutation({
        mutationFn: (payload: CreateRulesetRequest) => createRuleset(accessToken!, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ["admin", "rulesets"] });
            qc.invalidateQueries({ queryKey: ["admin", "rulesets", "current"] });
        },
    });
}
