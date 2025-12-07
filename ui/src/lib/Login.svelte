<script lang="ts">
  import { untrack } from "svelte";
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";

  interface ProvidersResponse {
    Providers: string[];
    ProvidersMap: Record<string, string>;
  }

  let providers = $state<ProvidersResponse | null>(null);
  let error = $state<string | null>(null);

  async function loadProviders() {
    error = null;
    try {
      const res = await fetch("/api/providers");
      if (!res.ok) {
        throw new Error("Failed to load providers");
      }
      providers = await res.json();
    } catch (e) {
      error = (e as Error).message;
    }
  }

  $effect(() => {
    untrack(() => {
      loadProviders();
    });
  });
</script>

<div class="flex min-h-screen items-center justify-center p-4">
  <Card.Root class="w-full max-w-md">
    <Card.Header>
      <Card.Title class="text-2xl font-bold text-center"
        >GH Secret Broker</Card.Title
      >
      <Card.Description class="text-center">
        Login with your GitHub account to manage secrets.
      </Card.Description>
    </Card.Header>
    <Card.Content class="grid gap-4">
      {#if error}
        <div class="text-destructive text-center text-sm">{error}</div>
        <Button variant="outline" onclick={() => loadProviders()}>Retry</Button>
      {:else if providers}
        <div class="grid gap-2">
          {#each providers.Providers as provider (provider)}
            <Button variant="outline" href="/auth/{provider}" class="w-full">
              Log in with {providers.ProvidersMap[provider]}
            </Button>
          {/each}
        </div>
      {:else}
        <div class="text-center text-muted-foreground">
          Loading providers...
        </div>
      {/if}
    </Card.Content>
  </Card.Root>
</div>
