import React, { useEffect, useMemo, useState } from "react";
import { humanizeError } from "@/shared/api/errors";
import { useCreateRuleset, useCurrentRuleset, useRulesets } from "@/entities/admin/queries";

function cx(...xs: Array<string | false | null | undefined>) {
    return xs.filter(Boolean).join(" ");
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

function SaveBtn(props: React.ButtonHTMLAttributes<HTMLButtonElement>) {
    return (
        <button
            {...props}
            className={cx(
                "h-[30px] px-10 rounded-[10px] border border-white/10",
                "bg-white/5 text-[12px] text-white/85",
                "hover:bg-white/10 disabled:opacity-50",
                props.className
            )}
        />
    );
}

function GhostBtn(props: React.ButtonHTMLAttributes<HTMLButtonElement>) {
    return (
        <button
            {...props}
            className={cx(
                "h-[30px] px-6 rounded-[10px] border border-white/10",
                "bg-transparent text-[12px] text-white/75",
                "hover:bg-white/5 disabled:opacity-50",
                props.className
            )}
        />
    );
}

function DangerMiniBtn(props: React.ButtonHTMLAttributes<HTMLButtonElement>) {
    return (
        <button
            {...props}
            className={cx(
                "h-[26px] px-4 rounded-[8px] border border-white/10",
                "text-[12px] text-red-300/90",
                "hover:bg-white/5 disabled:opacity-50 disabled:hover:bg-transparent",
                props.className
            )}
        />
    );
}

function MiniToggleBtn(props: React.ButtonHTMLAttributes<HTMLButtonElement> & { active?: boolean }) {
    const { active, ...rest } = props;
    return (
        <button
            {...rest}
            className={cx(
                "h-[26px] px-3 rounded-[8px] border border-white/10",
                "text-[12px] text-white/75",
                active ? "bg-white/5" : "bg-transparent",
                "hover:bg-white/5 disabled:opacity-50 disabled:hover:bg-transparent",
                props.className
            )}
        />
    );
}

type DraftLevel = {
    levelCode: string;
    thresholdTotalSpend: string;
    percentEarn: string;
};

const DEC_RE = /^\d+(\.\d{1,2})?$/;

function parseEffectiveFrom(v: string): string | null {
    const s = v.trim();
    if (!s || s.toLowerCase() === "now") return new Date().toISOString();
    const d = new Date(s);
    if (Number.isNaN(d.getTime())) return null;
    return d.toISOString();
}

function normalizeDec(v: string) {
    return v.replace(/\s+/g, "").replace(",", ".");
}

function fmtTs(iso?: string) {
    if (!iso) return "—";
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) return iso;

    const pad = (n: number) => String(n).padStart(2, "0");
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(
        d.getMinutes()
    )}`;
}

function validatePayload(input: { baseRubPerPoint: string; effectiveFrom: string; levels: DraftLevel[] }) {
    const base = normalizeDec(input.baseRubPerPoint);
    if (!base || !DEC_RE.test(base) || Number(base) <= 0) {
        return "baseRubPerPoint: invalid. Example: 10.00 (must be > 0)";
    }

    const iso = parseEffectiveFrom(input.effectiveFrom);
    if (!iso) return "effective_from: invalid date. Use ISO string or 'now'.";

    if (!input.levels.length) return "levels: at least one level is required.";

    for (let i = 0; i < input.levels.length; i++) {
        const row = input.levels[i];
        const lc = row.levelCode.trim();
        if (!lc) return `levels[${i}].levelCode: required`;

        const thr = normalizeDec(row.thresholdTotalSpend);
        if (!thr || !DEC_RE.test(thr) || Number(thr) < 0) {
            return `levels[${i}].thresholdTotalSpend: invalid. Example: 5000.00 (>= 0)`;
        }

        const pct = normalizeDec(row.percentEarn);
        if (!pct || !DEC_RE.test(pct) || Number(pct) <= 0) {
            return `levels[${i}].percentEarn: invalid. Example: 110.00 (> 0)`;
        }
    }

    return null;
}

export function AdminRulesPage() {
    const currentQ = useCurrentRuleset();
    const createM = useCreateRuleset();

    const current = (currentQ.data as any) ?? null;

    const [baseRubPerPoint, setBaseRubPerPoint] = useState("10.00");
    const [effectiveFrom, setEffectiveFrom] = useState("now");
    const [levelsDraft, setLevelsDraft] = useState<DraftLevel[]>([]);
    const [localError, setLocalError] = useState<string | null>(null);

    // Rulesets list
    const [offset, setOffset] = useState(0);
    const limit = 20;
    const rulesetsQ = useRulesets({ limit, offset });

    // раскрытие уровней в списке rulesets
    const [expandedIds, setExpandedIds] = useState<Set<number>>(() => new Set());
    const toggleExpanded = (id: number) => {
        setExpandedIds((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    };

    // подтягиваем current в форму
    useEffect(() => {
        if (!current) return;

        if (current.baseRubPerPoint) setBaseRubPerPoint(String(current.baseRubPerPoint));

        const srcLevels: any[] = Array.isArray(current.levels) ? current.levels : [];
        const sorted = [...srcLevels].sort((a, b) => {
            const aa = Number(normalizeDec(String(a.thresholdTotalSpend ?? "0")));
            const bb = Number(normalizeDec(String(b.thresholdTotalSpend ?? "0")));
            return aa - bb;
        });

        setLevelsDraft(
            sorted.map((l) => ({
                levelCode: String(l.levelCode ?? ""),
                thresholdTotalSpend: String(l.thresholdTotalSpend ?? "0.00"),
                percentEarn: String(l.percentEarn ?? "100.00"),
            }))
        );
    }, [current]);

    const levelsForTable = useMemo(() => levelsDraft, [levelsDraft]);

    function updateLevel(idx: number, patch: Partial<DraftLevel>) {
        setLevelsDraft((prev) => {
            const next = [...prev];
            next[idx] = { ...next[idx], ...patch };
            return next;
        });
    }

    function addLevel() {
        setLevelsDraft((prev) => [...prev, { levelCode: "New level", thresholdTotalSpend: "0.00", percentEarn: "100.00" }]);
    }

    function removeLevel(idx: number) {
        setLevelsDraft((prev) => prev.filter((_, i) => i !== idx));
    }

    function resetFromCurrent() {
        if (!current) return;

        const srcLevels: any[] = Array.isArray(current.levels) ? current.levels : [];
        const sorted = [...srcLevels].sort((a, b) => {
            const aa = Number(normalizeDec(String(a.thresholdTotalSpend ?? "0")));
            const bb = Number(normalizeDec(String(b.thresholdTotalSpend ?? "0")));
            return aa - bb;
        });

        setLevelsDraft(
            sorted.map((l) => ({
                levelCode: String(l.levelCode ?? ""),
                thresholdTotalSpend: String(l.thresholdTotalSpend ?? "0.00"),
                percentEarn: String(l.percentEarn ?? "100.00"),
            }))
        );
    }

    async function onSave() {
        setLocalError(null);

        const err = validatePayload({ baseRubPerPoint, effectiveFrom, levels: levelsDraft });
        if (err) {
            setLocalError(err);
            return;
        }

        const iso = parseEffectiveFrom(effectiveFrom)!;

        const payload: any = {
            baseRubPerPoint: normalizeDec(baseRubPerPoint),
            effectiveFrom: iso,
            levels: levelsDraft.map((l) => ({
                levelCode: l.levelCode.trim(),
                thresholdTotalSpend: normalizeDec(l.thresholdTotalSpend),
                percentEarn: normalizeDec(l.percentEarn),
            })),
        };

        try {
            await createM.mutateAsync(payload);
            setEffectiveFrom("now");
            setOffset(0);
            setExpandedIds(new Set()); // чтобы после изменения не держать раскрытия старых
            rulesetsQ.refetch();
            currentQ.refetch();
        } catch {
            // humanizeError(createM.error) ниже
        }
    }

    const canSave =
        !createM.isPending &&
        !currentQ.isLoading &&
        baseRubPerPoint.trim().length > 0 &&
        levelsDraft.length > 0;

    const rulesetItems: any[] = (rulesetsQ.data as any)?.items ?? [];
    const canPrev = offset > 0;
    const canNext = rulesetItems.length === limit;

    return (
        <div className="min-w-0">
            <div className="text-[18px] font-semibold">Rules</div>
            <div className="text-[11px] text-white/55 mt-1 mb-4">Create a new ruleset</div>

            {/* Create ruleset */}
            <PanelCard className="px-5 py-5 mb-6">
                <div className="grid grid-cols-[160px_160px_1fr] gap-10 items-start">
                    <div>
                        <div className="text-[11px] text-white/55 mb-2">baseRubPerPoint</div>
                        <AInput
                            value={baseRubPerPoint}
                            onChange={(e) => setBaseRubPerPoint(e.target.value)}
                            placeholder="10.00"
                            inputMode="decimal"
                        />
                    </div>

                    <div>
                        <div className="text-[11px] text-white/55 mb-2">effective_from</div>
                        <AInput value={effectiveFrom} onChange={(e) => setEffectiveFrom(e.target.value)} placeholder="now or ISO" />
                    </div>

                    <div className="flex justify-end gap-3 pt-[22px]">
                        <GhostBtn type="button" onClick={resetFromCurrent} disabled={!current}>
                            Reset levels
                        </GhostBtn>
                        <SaveBtn type="button" onClick={onSave} disabled={!canSave}>
                            {createM.isPending ? "Saving..." : "Save ruleset"}
                        </SaveBtn>
                    </div>
                </div>

                {(localError || createM.isError || currentQ.isError) && (
                    <div className="mt-4 text-[12px] text-red-200/90">
                        {localError ?? humanizeError((createM.error as any) ?? (currentQ.error as any))}
                    </div>
                )}
            </PanelCard>

            {/* Levels editor */}
            <div className="flex items-center justify-between mb-2">
                <div className="text-[12px] font-semibold text-white/85">Levels</div>
                <GhostBtn type="button" onClick={addLevel}>
                    Add level
                </GhostBtn>
            </div>

            <PanelCard className="overflow-hidden mb-6">
                <div className="px-5 py-4 border-b border-white/10">
                    <div className="grid grid-cols-[1.2fr_1.2fr_1fr] gap-6 text-[11px] text-white/55">
                        <div>level_code</div>
                        <div>threshold_total_spend</div>
                        <div>percent_earn</div>
                    </div>
                </div>

                {currentQ.isLoading && levelsForTable.length === 0 && (
                    <div className="px-5 py-5 text-[13px] text-white/55">Loading...</div>
                )}

                {!currentQ.isLoading && levelsForTable.length === 0 && (
                    <div className="px-5 py-5 text-[13px] text-white/55">No levels</div>
                )}

                {levelsForTable.length > 0 && (
                    <div className="divide-y divide-white/10">
                        {levelsForTable.map((l, idx) => (
                            <div key={`${l.levelCode}-${idx}`} className="px-5 py-3 grid grid-cols-[1.2fr_1.2fr_1fr] gap-6 items-center">
                                <AInput value={l.levelCode} onChange={(e) => updateLevel(idx, { levelCode: e.target.value })} />

                                <AInput
                                    value={l.thresholdTotalSpend}
                                    onChange={(e) => updateLevel(idx, { thresholdTotalSpend: e.target.value })}
                                    placeholder="5000.00"
                                    inputMode="decimal"
                                />

                                <div className="flex items-center gap-3">
                                    <AInput
                                        value={l.percentEarn}
                                        onChange={(e) => updateLevel(idx, { percentEarn: e.target.value })}
                                        placeholder="110.00"
                                        inputMode="decimal"
                                    />
                                    <DangerMiniBtn
                                        type="button"
                                        onClick={() => removeLevel(idx)}
                                        disabled={levelsDraft.length <= 1}
                                        title="Remove level"
                                    >
                                        Remove
                                    </DangerMiniBtn>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </PanelCard>

            {/* Rulesets list */}
            <div className="flex items-center justify-between mb-2">
                <div className="text-[12px] font-semibold text-white/85">Rulesets</div>
                <div className="flex gap-2">
                    <GhostBtn type="button" disabled={!canPrev} onClick={() => setOffset((x) => Math.max(0, x - limit))}>
                        Prev
                    </GhostBtn>
                    <GhostBtn type="button" disabled={!canNext} onClick={() => setOffset((x) => x + limit)}>
                        Next
                    </GhostBtn>
                </div>
            </div>

            <PanelCard className="overflow-hidden">
                <div className="px-5 py-4 border-b border-white/10">
                    <div className="grid grid-cols-[.6fr_1.4fr_1fr_.7fr] gap-6 text-[11px] text-white/55">
                        <div>id</div>
                        <div>effective_from</div>
                        <div>base</div>
                        <div>levels</div>
                    </div>
                </div>

                {rulesetsQ.isLoading && (
                    <div className="px-5 py-5 text-[13px] text-white/55">Loading...</div>
                )}

                {rulesetsQ.isError && (
                    <div className="px-5 py-5 text-[13px] text-red-200/90">
                        {humanizeError(rulesetsQ.error)}
                        <div className="mt-3">
                            <button
                                className="h-[30px] px-5 rounded-[8px] border border-white/10 text-[12px] hover:bg-white/5"
                                type="button"
                                onClick={() => rulesetsQ.refetch()}
                            >
                                Retry
                            </button>
                        </div>
                    </div>
                )}

                {!rulesetsQ.isLoading && !rulesetsQ.isError && rulesetItems.length === 0 && (
                    <div className="px-5 py-5 text-[13px] text-white/55">No rulesets</div>
                )}

                {!rulesetsQ.isLoading && !rulesetsQ.isError && rulesetItems.length > 0 && (
                    <div className="divide-y divide-white/10">
                        {rulesetItems.map((r: any) => {
                            const isCurrent = current?.id != null && r?.id === current.id;
                            const isOpen = expandedIds.has(Number(r.id));
                            const lvlCount = Array.isArray(r.levels) ? r.levels.length : 0;

                            return (
                                <React.Fragment key={r.id}>
                                    {/* main row */}
                                    <div className="px-5 py-4 grid grid-cols-[.6fr_1.4fr_1fr_.7fr] gap-6 items-center text-[12px] text-white/85">
                                        <div className="flex items-center gap-2">
                                            <span className="tabular-nums">{r.id}</span>
                                            {isCurrent && (
                                                <span className="text-[10px] px-2 py-1 rounded-[999px] border border-white/10 bg-white/5 text-white/70">
                          current
                        </span>
                                            )}
                                        </div>

                                        <div className="text-white/80">{fmtTs(r.effectiveFrom)}</div>
                                        <div className="tabular-nums text-white/80">{String(r.baseRubPerPoint)}</div>

                                        <div className="flex items-center justify-between gap-3">
                                            <span className="tabular-nums text-white/70">{lvlCount}</span>
                                            <MiniToggleBtn
                                                type="button"
                                                active={isOpen}
                                                onClick={() => toggleExpanded(Number(r.id))}
                                                disabled={!Array.isArray(r.levels) || lvlCount === 0}
                                                title={isOpen ? "Hide levels" : "Show levels"}
                                            >
                                                {isOpen ? "Hide" : "Show"}
                                            </MiniToggleBtn>
                                        </div>
                                    </div>

                                    {/* expanded levels */}
                                    {isOpen && Array.isArray(r.levels) && r.levels.length > 0 && (
                                        <div className="px-5 pb-4">
                                            <div className="rounded-[10px] border border-white/10 bg-[#0A1020] overflow-hidden">
                                                <div className="px-4 py-3 border-b border-white/10">
                                                    <div className="grid grid-cols-[1.2fr_1.2fr_1fr] gap-6 text-[11px] text-white/55">
                                                        <div>level_code</div>
                                                        <div>threshold_total_spend</div>
                                                        <div>percent_earn</div>
                                                    </div>
                                                </div>

                                                <div className="divide-y divide-white/10">
                                                    {r.levels.map((l: any) => (
                                                        <div
                                                            key={l.id ?? `${l.levelCode}-${l.thresholdTotalSpend}-${l.percentEarn}`}
                                                            className="px-4 py-3 grid grid-cols-[1.2fr_1.2fr_1fr] gap-6 items-center text-[12px] text-white/80"
                                                        >
                                                            <div className="text-white/85">{String(l.levelCode)}</div>
                                                            <div className="tabular-nums text-white/70">{String(l.thresholdTotalSpend)}</div>
                                                            <div className="tabular-nums text-white/70">{String(l.percentEarn)}</div>
                                                        </div>
                                                    ))}
                                                </div>
                                            </div>
                                        </div>
                                    )}
                                </React.Fragment>
                            );
                        })}
                    </div>
                )}
            </PanelCard>
        </div>
    );
}
