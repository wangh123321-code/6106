<script>
  import { api } from '$lib/api';
  import { onMount } from 'svelte';

  let medalBoard = [];
  let tournamentId = '';

  async function loadMedalBoard() {
    if (!tournamentId) return;
    const res = await api.get(`/tournaments/${tournamentId}/medal-board`);
    if (res.data) medalBoard = res.data;
  }
</script>

<svelte:head>
  <title>奖牌榜 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>奖牌榜</h1>
  <button class="btn btn-primary" on:click={loadMedalBoard}>查询</button>
</div>

<div class="card" style="margin-bottom:16px;">
  <div class="form-group" style="margin:0;">
    <label>赛事ID</label>
    <div style="display:flex;gap:8px;">
      <input type="text" bind:value={tournamentId} placeholder="输入赛事ID" />
    </div>
  </div>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>排名</th>
        <th>选手</th>
        <th>🥇 金牌</th>
        <th>🥈 银牌</th>
        <th>🥉 铜牌</th>
        <th>总计</th>
      </tr>
    </thead>
    <tbody>
      {#each medalBoard as entry, i}
        <tr>
          <td>{i + 1}</td>
          <td>{entry.display_name}</td>
          <td>{entry.gold}</td>
          <td>{entry.silver}</td>
          <td>{entry.bronze}</td>
          <td><strong>{entry.total}</strong></td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
