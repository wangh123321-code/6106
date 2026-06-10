<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  let tournament = null;
  let events = [];
  let user = null;

  $: id = $page.params.id;

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    user = getCurrentUser();
    await loadTournament();
  });

  async function loadTournament() {
    const res = await api.get(`/tournaments/${id}`);
    if (res.data) tournament = res.data;
    const evRes = await api.get(`/tournaments/${id}/events`);
    if (evRes.data) events = evRes.data;
  }

  function typeLabel(t) {
    const map = { ms: '男子单打', ws: '女子单打', xd: '混合双打' };
    return map[t] || t;
  }

  function drawMethodLabel(d) {
    return d === 'snake' ? '积分蛇形分组' : '随机抽签';
  }
</script>

<svelte:head>
  <title>{tournament?.name || '赛事详情'} - 赛事管理平台</title>
</svelte:head>

{#if tournament}
<div class="page-header">
  <h1>{tournament.name}</h1>
  <div style="display:flex;gap:8px;">
    {#if user?.role === 'player'}
      <a href="/tournaments/{id}/register" class="btn btn-primary">报名参赛</a>
    {/if}
    {#if user?.role === 'committee'}
      <a href="/tournaments/{id}/events/create" class="btn btn-primary">添加项目</a>
    {/if}
  </div>
</div>

<div class="card">
  <div class="info-grid">
    <div><span class="label">地点</span>{tournament.location}</div>
    <div><span class="label">时间</span>{tournament.start_date?.substring(0,10)} ~ {tournament.end_date?.substring(0,10)}</div>
    <div><span class="label">状态</span><span class="badge badge-info">{tournament.status}</span></div>
  </div>
</div>

<h2 style="margin:24px 0 12px;">比赛项目</h2>
<div class="events-grid">
  {#each events as e}
    <div class="card event-card">
      <h3>{e.name}</h3>
      <p class="text-secondary">{typeLabel(e.type)}</p>
      <p class="text-secondary">抽签方式：{drawMethodLabel(e.draw_method)}</p>
      <p class="text-secondary">
        状态：<span class="badge badge-info">{e.status}</span>
      </p>
      <div style="margin-top:12px;display:flex;gap:8px;">
        {#if e.status === 'drawn' || e.status === 'in_progress'}
          <a href="/events/{e.id}/bracket" class="btn btn-primary" style="padding:4px 12px;">查看对阵</a>
        {/if}
        {#if e.status === 'open' && user?.role === 'player'}
          <a href="/events/{e.id}/register" class="btn btn-success" style="padding:4px 12px;">报名</a>
        {/if}
        {#if e.registration_open && user?.role === 'committee'}
          <a href="/events/{e.id}/registrations" class="btn btn-outline" style="padding:4px 12px;">审核报名</a>
        {/if}
      </div>
    </div>
  {/each}
</div>
{/if}

<style>
  .info-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
  .label { color: var(--text-secondary); font-size: 13px; display: block; margin-bottom: 2px; }
  .events-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
  .event-card h3 { font-size: 16px; margin-bottom: 8px; }
  .text-secondary { color: var(--text-secondary); font-size: 14px; margin-bottom: 4px; }
</style>
