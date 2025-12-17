<script lang="ts">
    import { Input } from "$lib/components/ui/input";
    import { cn } from "$lib/utils.js";
    import type { HTMLInputAttributes } from "svelte/elements";

    let {
        class: className,
        value = $bindable(),
        ...rest
    }: HTMLInputAttributes = $props();
</script>

<Input class={cn("w-full", className)} bind:value {...rest as any} />
