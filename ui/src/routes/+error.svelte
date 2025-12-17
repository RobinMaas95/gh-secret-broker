<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";

    let { error, status } = $props<{ error: Error; status: number }>();

    const code = $derived(Number.isFinite(status) ? status : 500);
    const codeLabel = $derived(
        Number.isFinite(status) ? String(status) : "Error",
    );
    const isNotFound = $derived(code === 404);
    const title = $derived(
        isNotFound ? "Page not found" : "Something went wrong",
    );
    const message = $derived(
        isNotFound && !error?.message
            ? "We couldn't find the page you're looking for."
            : error?.message || "An unexpected error occurred.",
    );
</script>

<div class="min-h-[60vh] flex items-center justify-center px-4">
    <Card.Root class="w-full max-w-xl">
        <Card.Header>
            <Card.Title class="text-3xl font-bold">
                {codeLabel}: {title}
            </Card.Title>
            <Card.Description>{message}</Card.Description>
        </Card.Header>
        <Card.Footer class="gap-2">
            <Button href="/" variant="default">Go to Dashboard</Button>
            <Button href="/login" variant="outline">Login</Button>
        </Card.Footer>
    </Card.Root>
</div>
