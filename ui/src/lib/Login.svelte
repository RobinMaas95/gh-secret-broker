<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";

  import type { ProvidersResponse } from "$lib/types";

  let { providers } = $props<{ providers: ProvidersResponse | null }>();
</script>

<div class="flex min-h-[calc(100vh-60px)] items-center justify-center p-4">
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
      {#if providers}
        <div class="grid gap-2">
          {#each providers.Providers as provider (provider)}
            <Button variant="default" href="/auth/{provider}" class="w-full">
              Log in with {providers.ProvidersMap[provider]}
            </Button>
          {/each}
        </div>
      {:else}
        <div class="text-destructive text-center text-sm">
          Failed to load providers. Please try again later.
        </div>
      {/if}
    </Card.Content>
  </Card.Root>
</div>
