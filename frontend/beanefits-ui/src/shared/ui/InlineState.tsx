import React from "react";
import { Button } from "@/shared/ui/Button";
import { Callout } from "@/shared/ui/Callout";

export function LoadingBlock({ label = "Loading..." }: { label?: string }) {
    return (
        <div className="rounded-3xl border-2 border-border bg-white px-6 py-6 text-xl text-muted">
            {label}
        </div>
    );
}

export function ErrorBlock({
                               message,
                               onRetry,
                           }: {
    message: string;
    onRetry?: () => void;
}) {
    return (
        <div className="space-y-4">
            <Callout title="Something went wrong">{message}</Callout>
            {onRetry && (
                <Button variant="outline" fullWidth type="button" onClick={onRetry}>
                    Retry
                </Button>
            )}
        </div>
    );
}
