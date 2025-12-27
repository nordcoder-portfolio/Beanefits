import React from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const qc = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 10_000,
            retry: 1,
            refetchOnWindowFocus: false,
        },
    },
});

export function QueryProvider({ children }: { children: React.ReactNode }) {
    return <QueryClientProvider client={qc}>{children}</QueryClientProvider>;
}
