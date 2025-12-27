import React from "react";
import { NavLink } from "react-router-dom";

function Item({ to, label }: { to: string; label: string }) {
    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                [
                    "flex-1 py-6 text-center text-lg",
                    "relative",
                    isActive ? "text-primary" : "text-primary2",
                ].join(" ")
            }
        >
            {({ isActive }) => (
                <>
                    {isActive && (
                        <span className="absolute top-1 left-1/2 -translate-x-1/2 h-2 w-24 rounded-full bg-primary" />
                    )}
                    <span className="block mt-1">{label}</span>
                </>
            )}
        </NavLink>
    );
}

export function BottomNav() {
    return (
        <div className="flex items-center">
            <Item to="/home" label="Home" />
            <Item to="/history" label="History" />
            <Item to="/profile" label="Profile" />
        </div>
    );
}
