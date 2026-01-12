import { requestJson } from "@/shared/api/client";
import type {
    CreateRulesetRequest,
    Ruleset,
    RulesetsPage,
    UsersPage,
} from "@/shared/api/contracts";

export async function listUsers(
    token: string,
    params: { limit?: number; offset?: number; q?: string } = {}
): Promise<UsersPage> {
    const qs = new URLSearchParams();
    qs.set("limit", String(params.limit ?? 20));
    qs.set("offset", String(params.offset ?? 0));
    if (params.q) qs.set("q", params.q);

    return requestJson<UsersPage>(`/admin/users?${qs.toString()}`, { token });
}

export async function deleteUser(token: string, userId: number): Promise<void> {
    await requestJson<void>(`/admin/users/${userId}`, { method: "DELETE", token });
}

export async function getCurrentRuleset(token: string): Promise<Ruleset> {
    return requestJson<Ruleset>("/admin/rulesets/current", { token });
}

export async function listRulesets(
    token: string,
    params: { limit?: number; offset?: number } = {}
): Promise<RulesetsPage> {
    const qs = new URLSearchParams();
    qs.set("limit", String(params.limit ?? 20));
    qs.set("offset", String(params.offset ?? 0));

    return requestJson<RulesetsPage>(`/admin/rulesets?${qs.toString()}`, { token });
}

export async function createRuleset(
    token: string,
    payload: CreateRulesetRequest
): Promise<Ruleset> {
    return requestJson<Ruleset>("/admin/rulesets", {
        method: "POST",
        token,
        body: JSON.stringify(payload),
    });
}

const api = { listUsers, deleteUser, getCurrentRuleset, listRulesets, createRuleset };
export default api;
