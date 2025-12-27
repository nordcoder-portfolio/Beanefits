// src/shared/fixtures/rulesets.ts

import type { Ruleset, RulesetsPage, LevelRule } from "@/shared/api/contracts";
import { isoDaysAgo } from "./clock";

const LEVELS: Omit<LevelRule, "id">[] = [
    { levelCode: "Green Bean",   thresholdTotalSpend: "0.00",     percentEarn: "100.00" },
    { levelCode: "Light Roast",  thresholdTotalSpend: "3000.00",  percentEarn: "105.00" },
    { levelCode: "Medium Roast", thresholdTotalSpend: "7000.00",  percentEarn: "110.00" },
    { levelCode: "Dark Roast",   thresholdTotalSpend: "15000.00", percentEarn: "115.00" },
    { levelCode: "Premium Roast",thresholdTotalSpend: "30000.00", percentEarn: "120.00" },
];

export const fixtureCurrentRuleset: Ruleset = {
    id: 200,
    effectiveFrom: isoDaysAgo(14),
    baseRubPerPoint: "10.00",
    levels: LEVELS.map((l, i) => ({ id: 2000 + i, ...l })),
    createdAt: isoDaysAgo(14),
};

export const fixtureOldRuleset: Ruleset = {
    id: 199,
    effectiveFrom: isoDaysAgo(60),
    baseRubPerPoint: "12.00",
    levels: LEVELS.map((l, i) => ({
        id: 1900 + i,
        ...l,
        percentEarn: i === 0 ? "100.00" : String(100 + i * 3).padEnd(6, "0"), // чуть иные проценты
    })),
    createdAt: isoDaysAgo(60),
};

export function makeRulesetsPage(): RulesetsPage {
    return {
        items: [fixtureCurrentRuleset, fixtureOldRuleset],
        total: 2,
    };
}
