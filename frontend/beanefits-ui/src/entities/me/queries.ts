import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import { useAuth } from "@/app/providers/AuthProvider";
import * as meApi from "./api.ts";

export function useMeQuery() {
    const { accessToken } = useAuth();
    return useQuery({
        queryKey: ["me"],
        queryFn: () => meApi.getMe(accessToken!),
        enabled: Boolean(accessToken),
    });
}

export function useBalanceQuery() {
    const { accessToken } = useAuth();
    return useQuery({
        queryKey: ["me", "balance"],
        queryFn: () => meApi.getMyBalance(accessToken!),
        enabled: Boolean(accessToken),
    });
}

export function useEventsInfinite(limit = 20) {
    const { accessToken } = useAuth();
    return useInfiniteQuery({
        queryKey: ["me", "events", limit],
        enabled: Boolean(accessToken),
        queryFn: ({ pageParam }) => meApi.getMyEvents(accessToken!, limit, pageParam ?? null),
        initialPageParam: null as string | null,
        getNextPageParam: (last) => last.nextBeforeTs ?? undefined,
    });
}
