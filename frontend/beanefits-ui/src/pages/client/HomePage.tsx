import React from "react";
import { QRCodeCanvas } from "qrcode.react";
import { Card } from "@/shared/ui/Card";
import { Button } from "@/shared/ui/Button";
import { copyToClipboard } from "@/shared/lib/clipboard";
import { useBalanceQuery, useMeQuery } from "@/entities/me/queries";
import { humanizeError } from "@/shared/api/errors";
import { ErrorBlock, LoadingBlock } from "@/shared/ui/InlineState";

export function HomePage() {
    const me = useMeQuery();
    const bal = useBalanceQuery();

    if (me.isLoading || bal.isLoading) {
        return (
            <div className="px-7 pt-10">
                <LoadingBlock label="Loading your account..." />
            </div>
        );
    }

    if (me.isError || bal.isError) {
        const msg = humanizeError((me.error as any) ?? (bal.error as any));
        return (
            <div className="px-7 pt-10">
                <ErrorBlock message={msg} onRetry={() => { me.refetch(); bal.refetch(); }} />
            </div>
        );
    }

    const publicCode = me.data?.account.publicCode ?? "ABCD-1234-XYZ";
    const points = bal.data?.balancePoints ?? me.data?.account.balancePoints ?? 0;
    const level = bal.data?.levelCode ?? me.data?.account.levelCode ?? "Green Bean";

    return (
        <div className="px-7 pt-10">
            <Card className="p-8">
                <div className="text-xl text-muted mb-6">Show at checkout</div>

                {/* Уровень над QR */}
                <div className="mb-5 flex items-center justify-between">
                    <div className="text-xl text-muted">Loyalty level</div>
                    <div className="rounded-full border-2 border-[#F3D3A8] bg-peach px-8 py-4 text-xl text-primary2 whitespace-nowrap">
                        {level}
                    </div>
                </div>

                {/* QR */}
                <div className="rounded-3xl border-2 border-[#F3D3A8] bg-peach p-6 flex items-center justify-center">
                    <QRCodeCanvas value={publicCode} size={260} />
                </div>

                <div className="text-center text-primary2 mt-5 text-xl">QR</div>

                {/* Points без перекрытий */}
                <div className="mt-8 rounded-3xl border-2 border-accent bg-[#FFF9E6] px-6 py-5">
                    <div className="text-xl text-muted">Points</div>
                    <div className="text-5xl font-medium tracking-wide mt-1 tabular-nums overflow-hidden text-ellipsis whitespace-nowrap">
                        {points.toLocaleString("en-US")}
                    </div>
                </div>

                <div className="mt-10 text-xl">Public code</div>

                <div className="mt-4 flex gap-4 items-stretch">
                    <div className="flex-1 rounded-3xl border-2 border-[#F3D3A8] bg-peach px-6 py-5 text-xl text-primary2 flex items-center overflow-hidden">
                        <span className="truncate">{publicCode}</span>
                    </div>
                    <Button
                        type="button"
                        className="px-10"
                        onClick={async () => copyToClipboard(publicCode)}
                    >
                        Copy
                    </Button>
                </div>
            </Card>
        </div>
    );
}
