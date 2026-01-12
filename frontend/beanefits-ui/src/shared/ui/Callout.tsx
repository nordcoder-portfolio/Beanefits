import React from "react";

export function Callout({
                            title = "Error",
                            children,
                        }: {
    title?: string;
    children: React.ReactNode;
}) {
    return (
        <div className="rounded-3xl border-2 border-danger bg-white px-6 py-5">
            <div className="text-xl font-medium text-danger">{title}</div>
            <div className="text-lg text-ink mt-2">{children}</div>
        </div>
    );
}
