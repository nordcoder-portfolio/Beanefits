import React from "react";

export function Input(
    props: React.InputHTMLAttributes<HTMLInputElement> & { hasError?: boolean }
) {
    const { hasError, className = "", ...rest } = props;
    return (
        <input
            className={[
                "w-full rounded-2xl bg-peach px-6 py-5 text-lg outline-none",
                "border-2",
                hasError ? "border-danger" : "border-accent",
                className,
            ].join(" ")}
            {...rest}
        />
    );
}
