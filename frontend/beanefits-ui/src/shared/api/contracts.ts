// src/shared/api/contracts.ts

export type RoleCode = "CLIENT" | "CASHIER" | "ADMIN";
export type EventType = "EARN" | "SPEND";
export type OperationType = "EARN" | "SPEND";

export interface Problem {
    type: string; // "about:blank"
    title: string;
    status: number;
    detail?: string;
    instance?: string;
    code?: string; // e.g. "NOT_ENOUGH_BALANCE"
}

export interface User {
    id: number;
    phone: string; // E.164-like
    roles: RoleCode[];
    isActive: boolean;
    createdAt: string; // ISO
}

export interface Account {
    id: number;
    publicCode: string; // uuid
    balancePoints: number;
    totalSpendMoney: string; // decimal string
    levelCode: string; // free-form label
    createdAt: string; // ISO
}

export interface AuthResponse {
    accessToken: string;
    user: User;
    account: Account;
}

export interface ClientProfile {
    user: User;
    account: Account;
}

export interface BalanceResponse {
    accountId: number;
    balancePoints: number;
    totalSpendMoney: string;
    levelCode: string;
    asOf: string; // ISO
}

export interface Event {
    id: number;
    accountId: number;
    type: EventType;
    deltaPoints: number; // signed
    balanceAfter: number;
    amountMoney?: string | null; // present for EARN
    rulesetId?: number | null;
    actorUserId?: number | null;
    ts: string; // ISO
}

export interface EventsPage {
    items: Event[];
    nextBeforeTs?: string | null;
}

export interface UsersPage {
    items: User[];
    total?: number | null;
}

export interface LevelRule {
    id: number;
    levelCode: string;
    thresholdTotalSpend: string; // decimal string
    percentEarn: string; // decimal string
}

export interface Ruleset {
    id: number;
    effectiveFrom: string; // ISO
    baseRubPerPoint: string; // decimal string
    levels: LevelRule[];
    createdAt: string; // ISO
}

export interface RulesetsPage {
    items: Ruleset[];
    total?: number | null;
}

export interface OperationResult {
    operationId: string;
    opType: OperationType;
    event: Event;
    balance: BalanceResponse;
    idempotentReplay?: boolean;
}
