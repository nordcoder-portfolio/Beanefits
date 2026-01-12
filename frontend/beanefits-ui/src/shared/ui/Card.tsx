import React from "react";

export function Card({ className = "", ...props }: React.HTMLAttributes<HTMLDivElement>) {
    return <div className={`rounded-xl3 bg-card shadow-soft ${className}`} {...props} />;
}
