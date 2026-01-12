export type RoleCode = "CLIENT" | "CASHIER" | "ADMIN";
export type EventType = "EARN" | "SPEND";
export type OperationType = "EARN" | "SPEND";

export interface Problem {
    type: string;
    title: string;
    status: number;
    detail?: string;
    instance?: string;
    code?: string;
}

export interface User {
    id: number;
    phone: string;
    roles: RoleCode[];
    isActive: boolean;
    createdAt: string;
}

export interface Account {
    id: number;
    publicCode: string;
    balancePoints: number;
    totalSpendMoney: string;
    levelCode: string;
    createdAt: string;
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
    asOf: string;
}

export interface Event {
    id: number;
    accountId: number;
    type: EventType;
    deltaPoints: number;
    balanceAfter: number;
    amountMoney?: string | null;
    rulesetId?: number | null;
    actorUserId?: number | null;
    ts: string;
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
    thresholdTotalSpend: string;
    percentEarn: string;
}

export interface Ruleset {
    id: number;
    effectiveFrom: string;
    baseRubPerPoint: string;
    levels: LevelRule[];
    createdAt: string;
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
