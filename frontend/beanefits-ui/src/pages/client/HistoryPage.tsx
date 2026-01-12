import React from "react";
import { Card } from "@/shared/ui/Card";
import { formatTs } from "@/shared/lib/format";
import { useEventsInfinite } from "@/entities/me/queries";
import { Button } from "@/shared/ui/Button";
import { humanizeError } from "@/shared/api/errors";
import { ErrorBlock, LoadingBlock } from "@/shared/ui/InlineState";

export function HistoryPage() {
    const q = useEventsInfinite(20);

    if (q.isLoading) {
        return (
            <div className="px-7 pt-10">
                <LoadingBlock label="Loading history..." />
            </div>
        );
    }

    if (q.isError) {
        return (
            <div className="px-7 pt-10">
                <ErrorBlock message={humanizeError(q.error)} onRetry={() => q.refetch()} />
            </div>
        );
    }

    const items = q.data?.pages.flatMap((p) => p.items) ?? [];

    return (
        <div className="px-7 pt-10">
            <Card className="p-8">
                <div className="text-4xl font-medium mb-8">History</div>

                <div className="space-y-5">
                    {items.map((e, idx) => {
                        const isEarn = e.type === "EARN";
                        const border = idx % 2 === 0 ? "border-accent bg-[#FFF9E6]" : "border-border bg-white";
                        const deltaColor = isEarn ? "text-success" : "text-danger";
                        const delta = isEarn ? `+${e.deltaPoints}` : `${e.deltaPoints}`;

                        return (
                            <div key={e.id} className={`rounded-3xl border-2 ${border} p-6`}>
                                <div className="text-xl text-muted">{formatTs(e.ts)}</div>
                                <div className="mt-3 flex items-start justify-between">
                                    <div>
                                        <div className="text-2xl font-medium">{e.type}</div>
                                        <div className="text-xl text-muted mt-1">Balance: {e.balanceAfter}</div>
                                    </div>
                                    <div className={`text-2xl font-medium ${deltaColor}`}>{delta}</div>
                                </div>
                            </div>
                        );
                    })}

                    {q.hasNextPage && (
                        <Button
                            variant="outline"
                            fullWidth
                            disabled={q.isFetchingNextPage}
                            onClick={() => q.fetchNextPage()}
                            type="button"
                            className="mt-2"
                        >
                            Load more
                        </Button>
                    )}
                </div>
            </Card>
        </div>
    );
}
