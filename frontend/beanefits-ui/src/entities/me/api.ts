import type { BalanceResponse, ClientProfile, EventsPage } from "@/shared/api/contracts";
import { requestJson } from "@/shared/api/client";

export function getMe(token: string): Promise<ClientProfile> {
    return requestJson<ClientProfile>("/me", { token });
}

export function getMyBalance(token: string): Promise<BalanceResponse> {
    return requestJson<BalanceResponse>("/me/balance", { token });
}

export function getMyEvents(token: string, limit = 20, beforeTs?: string | null): Promise<EventsPage> {
    const qs = new URLSearchParams();
    qs.set("limit", String(limit));
    if (beforeTs) qs.set("beforeTs", beforeTs);
    return requestJson<EventsPage>(`/me/events?${qs.toString()}`, { token });
}
