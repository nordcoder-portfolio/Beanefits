import React from "react";
import { NavLink, Outlet } from "react-router-dom";

function cx(...xs: Array<string | false | null | undefined>) {
    return xs.filter(Boolean).join(" ");
}

function SideItem({ to, label }: { to: string; label: string }) {
    return (
        <NavLink
            to={to}
            end
            className={({ isActive }) =>
                cx(
                    "block w-full rounded-[10px] px-4 py-3 text-[14px] leading-none",
                    "border border-white/10",
                    isActive ? "bg-white/5 text-white" : "bg-transparent text-white/85 hover:bg-white/5"
                )
            }
        >
            {label}
        </NavLink>
    );
}

export function AdminLayout() {
    return (
        <div className="min-h-dvh bg-[#070E1B] text-white flex">
            {/* Sidebar */}
            <aside className="w-[220px] shrink-0 bg-gradient-to-b from-[#0C1528] to-[#081021] border-r border-white/5 px-6 py-8">
                <div className="mb-10">
                    <div className="text-[16px] font-semibold">Admin</div>
                    <div className="text-[11px] text-white/55 mt-1">Loyalty system</div>
                </div>

                <nav className="space-y-3">
                    <SideItem to="/admin/users" label="Users" />
                    <SideItem to="/admin/rules" label="Rules" />
                </nav>
            </aside>

            {/* Content */}
            <main className="flex-1 min-w-0 px-10 py-8">
                <Outlet />
            </main>
        </div>
    );
}
