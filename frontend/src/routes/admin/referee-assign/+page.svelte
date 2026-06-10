<script>
  import { api } from '$lib/api';
  import { onMount } from 'svelte';

  let tournaments = [];
  let events = [];
  let matches = [];
  let selectedEvent = '';
  let selectedMatch = '';
  let refereeId = '';
  let message = '';

  onMount(async () => {
    const res = await api.get('/tournaments');
    if (res.data) tournaments = res.data;
  });

  async function loadEvents(tournamentId) {
    const res = await api.get(`/tournaments/${tournamentId}/events`);
    if (res.data) events = res.data;
  }

  async function loadMatches(eventId) {
    const res = await api.get(`/events/${eventId}/bracket`);
    if (res.data?.matches) matches = res.data.matches;
  }

  async function assignReferee() {
    if (!selectedMatch || !refereeId) return;
    const res = await api.post(`/matches/${selectedMatch}/assign-referee`, {
      referee_id: refereeId
    });
    message = res.code === 0 ? '指派成功' : (res.message || '指派失败');
  }
</script>

<svelte:head>
  <title>裁判指派 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>裁判指派</h1>
</div>

<div class="card" style="max-width:600px;">
  {#if message}
    <div class="alert">{message}</div>
  {/if}

  <div class="form-group">
    <label>选择赛事</label>
    <select onchange={(e) => loadEvents(e.target.value)}>
      <option value="">请选择</option>
      {#each tournaments as t}
        <option value={t.id}>{t.name}</option>
      {/each}
    </select>
  </div>

  <div class="form-group">
    <label>选择项目</label>
    <select bind:value={selectedEvent} onchange={() => loadMatches(selectedEvent)}>
      <option value="">请选择</option>
      {#each events as e}
        <option value={e.id}>{e.name}</option>
      {/each}
    </select>
  </div>

  <div class="form-group">
    <label>选择场次</label>
    <select bind:value={selectedMatch}>
      <option value="">请选择</option>
      {#each matches as m}
        <option value={m.id}>第{m.round}轮-{m.position} ({m.player1_name || '待定'} vs {m.player2_name || '待定'})</option>
      {/each}
    </select>
  </div>

  <div class="form-group">
    <label>裁判用户ID</label>
    <input type="text" bind:value={refereeId} placeholder="输入裁判用户ID" />
  </div>

  <button class="btn btn-primary" style="width:100%;justify-content:center;" on:click={assignReferee}>
    指派裁判
  </button>
</div>

<style>
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; background: #e8f0fe; color: #1a73e8; }
</style>
