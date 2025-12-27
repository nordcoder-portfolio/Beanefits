// src/shared/fixtures/events.ts

import type { BalanceResponse, Event, EventsPage } from "@/shared/api/contracts";
import { isoMinutesAgo } from "./clock";

export function buildClientEvents(accountId: number): Event[] {
    // newest first
    // баланс всегда неотрицательный; SPEND — отрицательная deltaPoints
    const items: Event[] = [
        {
            id: 501,
            accountId,
            type: "EARN",
            deltaPoints: +30,
            balanceAfter: 240,
            amountMoney: "350.00",
            rulesetId: 200,
            actorUserId: 2,
            ts: isoMinutesAgo(20),
        },
        {
            id: 500,
            accountId,
            type: "SPEND",
            deltaPoints: -60,
            balanceAfter: 210,
            amountMoney: null,
            rulesetId: 200,
            actorUserId: 2,
            ts: isoMinutesAgo(60 * 6),
        },
        {
            id: 499,
            accountId,
            type: "EARN",
            deltaPoints: +120,
            balanceAfter: 270,
            amountMoney: "1400.00",
            rulesetId: 200,
            actorUserId: 2,
            ts: isoMinutesAgo(60 * 24),
        },
        {
            id: 498,
            accountId,
            type: "EARN",
            deltaPoints: +150,
            balanceAfter: 150,
            amountMoney: "1800.00",
            rulesetId: 200,
            actorUserId: 2,
            ts: isoMinutesAgo(60 * 48),
        },
    ];

    return items;
}

export function makeEventsPage(accountId: number, limit = 20): EventsPage {
    const items = buildClientEvents(accountId).slice(0, limit);
    const nextBeforeTs = items.length > 0 ? items[items.length - 1].ts : null;
    return { items, nextBeforeTs };
}

export function makeBalanceFromLatestEvent(accountId: number, events: Event[]): BalanceResponse {
    const latest = events[0];
    return {
        accountId,
        balancePoints: latest ? latest.balanceAfter : 0,
        totalSpendMoney: "7800.00",
        levelCode: "Light Roast",
        asOf: latest ? latest.ts : new Date().toISOString(),
    };
}
