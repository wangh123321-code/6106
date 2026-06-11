<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import BracketTree from '$lib/components/BracketTree.svelte';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let bracket = null;
  let matches = [];
  let event = null;
  let user = null;

  $: eventId = $page.params.event_id;

  onMount(async () => {
    if (!isLoggedIn()) return;
    user = getCurrentUser();
    await loadBracket();
  });

  async function loadBracket() {
    const res = await api.get(`/events/${eventId}/bracket`);
    if (res.data) {
      bracket = res.data.bracket;
      matches = res.data.matches || [];
    }
  }

  function refresh() {
    loadBracket();
  }
</script>

<svelte:head>
  <title>对阵表 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>对阵表</h1>
  <button class="btn btn-outline" on:click={refresh}>刷新</button>
</div>

{#if matches.length > 0}
  <BracketTree {matches} {event} />
{:else}
  <div class="card" style="text-align:center;padding:48px;">
    <p style="color:var(--text-secondary);">暂无对阵数据，请等待组委会进行抽签</p>
  </div>
{/if}
