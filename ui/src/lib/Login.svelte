<script lang="ts">
  import { untrack } from 'svelte';

  interface ProvidersResponse {
    Providers: string[];
    ProvidersMap: Record<string, string>;
  }

  let providers = $state<ProvidersResponse | null>(null);
  let error = $state<string | null>(null);

  async function loadProviders() {
    error = null;
    try {
      const res = await fetch('/api/providers');
      if (!res.ok) {
        throw new Error('Failed to load providers');
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

<div class="container">
  <h1>GH Secret Broker</h1>
  <div class="card">
    {#if error}
      <p class="error">{error}</p>
      <button onclick={() => loadProviders()}>Retry</button>
    {:else if providers}
      {#each providers.Providers as provider (provider)}
        <p><a href="/auth/{provider}">Log in with {providers.ProvidersMap[provider]}</a></p>
      {/each}
    {:else}
      <p>Loading providers...</p>
    {/if}
  </div>
</div>

<style>
  .container {
    max-width: 1280px;
    margin: 0 auto;
    padding: 2rem;
    text-align: center;
  }
  h1 {
    font-size: 3.2em;
    line-height: 1.1;
  }
  .card {
    padding: 2em;
    background-color: #1a1a1a;
    border-radius: 8px;
    margin-top: 1em;
  }
  .error {
    color: red;
  }
</style>
