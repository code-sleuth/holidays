<script lang="ts">
    import { cn } from "$lib/utils";
    import type { Snippet } from "svelte";
    import type { HTMLButtonAttributes } from "svelte/elements";

    interface Props extends HTMLButtonAttributes {
        variant?: "default" | "destructive" | "outline" | "secondary" | "ghost";
        size?: "default" | "sm" | "lg" | "icon";
        children: Snippet;
    }

    let {
        variant = "default",
        size = "default",
        class: className,
        children,
        ...rest
    }: Props = $props();

    const variants: Record<string, string> = {
        default: "bg-slate-900 text-white hover:bg-slate-800",
        destructive: "bg-red-500 text-white hover:bg-red-600",
        outline: "border border-slate-200 bg-white hover:bg-slate-100",
        secondary: "bg-slate-100 text-slate-900 hover:bg-slate-200",
        ghost: "hover:bg-slate-100",
    };

    const sizes: Record<string, string> = {
        default: "h-10 px-4 py-2",
        sm: "h-9 px-3 text-sm",
        lg: "h-11 px-8",
        icon: "h-10 w-10",
    };
</script>

<button
    class={cn(
        "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 disabled:pointer-events-none disabled:opacity-50 cursor-pointer",
        variants[variant],
        sizes[size],
        className,
    )}
    {...rest}
>
    {@render children()}
</button>
