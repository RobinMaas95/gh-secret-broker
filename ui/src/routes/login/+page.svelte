<script lang="ts">
    import { onMount } from "svelte";
    import Login from "$lib/Login.svelte";
    import { toast } from "svelte-sonner";
    import type { ProvidersResponse } from "$lib/types";

    let { data } = $props<{
        data: { providers: ProvidersResponse; unauthorized: boolean };
    }>();

    let unauthorizedToastShown = $state(false);

    function showUnauthorizedToastOnce() {
        if (!unauthorizedToastShown) {
            unauthorizedToastShown = true;
            queueMicrotask(() => toast.error("Please log in to continue"));
        }
    }

    $effect(() => {
        if (data.unauthorized) {
            showUnauthorizedToastOnce();
        }
    });

    onMount(() => {
        const unauthorized = new URLSearchParams(window.location.search).has(
            "unauthorized",
        );
        if (unauthorized) {
            showUnauthorizedToastOnce();
        }
    });
</script>

<Login providers={data.providers} />
