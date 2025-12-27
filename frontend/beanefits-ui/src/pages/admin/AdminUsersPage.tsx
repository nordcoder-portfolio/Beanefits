import React, { useEffect, useState } from "react";
import { humanizeError } from "@/shared/api/errors";
import { useAdminUsers, useDeleteUser } from "@/entities/admin/queries";

function cx(...xs: Array<string | false | null | undefined>) {
    return xs.filter(Boolean).join(" ");
}

function useDebouncedValue<T>(value: T, delayMs: number) {
    const [v, setV] = useState(value);
    useEffect(() => {
        const t = window.setTimeout(() => setV(value), delayMs);
        return () => window.clearTimeout(t);
    }, [value, delayMs]);
    return v;
}

function PanelCard({ children, className }: { children: React.ReactNode; className?: string }) {
    return (
        <div
            className={cx(
                "rounded-[12px] border border-white/10",
                "bg-gradient-to-b from-[#0E1A33] to-[#0B152B]",
                "shadow-[0_0_0_1px_rgba(255,255,255,0.03),0_30px_90px_rgba(0,0,0,0.45)]",
                className
            )}
        >
            {children}
        </div>
    );
}

function AInput(props: React.InputHTMLAttributes<HTMLInputElement>) {
    return (
        <input
            {...props}
            className={cx(
                "w-full h-[34px] rounded-[8px] border border-white/10 bg-[#0A1020]",
                "px-3 text-[13px] text-white/90 placeholder:text-white/35 outline-none",
                "focus:ring-2 focus:ring-white/10",
                props.className
            )}
        />
    );
}

function DeleteBtn(props: React.ButtonHTMLAttributes<HTMLButtonElement>) {
    return (
        <button
            {...props}
            className={cx(
                "h-[26px] px-6 rounded-[8px] border border-white/10",
                "text-[12px] text-red-300/90",
                "hover:bg-white/5 disabled:opacity-50 disabled:hover:bg-transparent"
            )}
        />
    );
}

export function AdminUsersPage() {
    const [q, setQ] = useState("");
    const dq = useDebouncedValue(q, 250);

    const limit = 20;
    const offset = 0;

    const usersQ = useAdminUsers({ limit, offset, q: dq.trim() ? dq.trim() : undefined });
    const del = useDeleteUser();

    const items = (usersQ.data as any)?.items ?? [];

    return (
        <div className="min-w-0">
            <div className="text-[18px] font-semibold mb-4">Users</div>

            {/* Search card */}
            <PanelCard className="px-5 py-4 mb-4">
                <div className="flex items-center gap-4">
                    <div className="text-[11px] text-white/55 w-[56px]">Search</div>
                    <div className="w-[260px]">
                        <AInput value={q} onChange={(e) => setQ(e.target.value)} placeholder="phone / role" />
                    </div>
                </div>
            </PanelCard>

            {/* Table card */}
            <PanelCard className="overflow-hidden">
                {/* header row */}
                <div className="px-5 py-4 border-b border-white/10">
                    <div className="grid grid-cols-[1.6fr_1.1fr_.6fr_.7fr] gap-6 text-[11px] text-white/55">
                        <div>Phone</div>
                        <div>Roles</div>
                        <div>Active</div>
                        <div>Action</div>
                    </div>
                </div>

                {/* body */}
                {usersQ.isLoading && (
                    <div className="px-5 py-5 text-[13px] text-white/55">Loading...</div>
                )}

                {usersQ.isError && (
                    <div className="px-5 py-5 text-[13px] text-red-200/90">
                        {humanizeError(usersQ.error)}
                        <div className="mt-3">
                            <button
                                className="h-[30px] px-5 rounded-[8px] border border-white/10 text-[12px] hover:bg-white/5"
                                type="button"
                                onClick={() => usersQ.refetch()}
                            >
                                Retry
                            </button>
                        </div>
                    </div>
                )}

                {!usersQ.isLoading && !usersQ.isError && items.length === 0 && (
                    <div className="px-5 py-5 text-[13px] text-white/55">No users</div>
                )}

                {!usersQ.isLoading && !usersQ.isError && items.length > 0 && (
                    <div className="divide-y divide-white/10">
                        {items.map((u: any) => (
                            <div
                                key={u.id}
                                className="px-5 py-4 grid grid-cols-[1.6fr_1.1fr_.6fr_.7fr] gap-6 items-center text-[12px] text-white/85"
                            >
                                <div>{u.phone}</div>
                                <div>{Array.isArray(u.roles) ? u.roles.join(", ") : String(u.roles ?? "")}</div>
                                <div className="text-white/70">{String(Boolean(u.isActive)).toLowerCase()}</div>
                                <div>
                                    <DeleteBtn
                                        type="button"
                                        disabled={del.isPending}
                                        onClick={() => del.mutate(u.id)}
                                    >
                                        Delete
                                    </DeleteBtn>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </PanelCard>
        </div>
    );
}
