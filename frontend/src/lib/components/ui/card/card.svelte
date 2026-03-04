<script lang="ts">
    import { cn } from "$lib/utils";
    import type { Snippet } from "svelte";
    import type { HTMLAttributes } from "svelte/elements";

    interface Props extends HTMLAttributes<HTMLDivElement> {
        children: Snippet;
    }

    let { class: className, children, ...rest }: Props = $props();
</script>

<!-- rounded-lg border border-slate-200 bg-white text-slate-950 shadow-sm -->
<div
    class={cn(
        "bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm",
        className,
    )}
    {...rest}
>
    {@render children?.()}
</div>
