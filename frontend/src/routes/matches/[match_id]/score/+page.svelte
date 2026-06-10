<script>
  import { api, getCurrentUser } from '$lib/api';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let match = null;
  let games = [];
  let bestOf = 7;
  let user = null;
  let saving = false;
  let error = '';

  $: matchId = $page.params.match_id;

  onMount(async () => {
    user = getCurrentUser();
    const res = await api.get(`/matches/${matchId}`);
    if (res.data) {
      match = res.data;
      bestOf = match.best_of || 7;
      if (match.games && match.games.length > 0) {
        games = match.games.map(g => ({ ...g }));
      } else {
        games = Array.from({ length: bestOf }, () => ({ score1: 0, score2: 0 }));
      }
    }
  });

  async function submitScore() {
    saving = true;
    error = '';
    const validGames = games.filter(g => g.score1 > 0 || g.score2 > 0);
    const res = await api.post(`/matches/${matchId}/score`, { games: validGames });
    if (res.code === 0) {
      goto('/referee');
    } else {
      error = res.message || '提交失败';
    }
    saving = false;
  }

  async function declareWalkover(playerId) {
    const reason = prompt('请输入弃权原因');
    if (!reason) return;
    const res = await api.post(`/matches/${matchId}/walkover`, {
      player_id: playerId,
      reason
    });
    if (res.code === 0) {
      goto('/referee');
    }
  }
</script>

<svelte:head>
  <title>录入比分 - 赛事管理平台</title>
</svelte:head>

{#if match}
<div class="page-header">
  <h1>录入比分</h1>
</div>

<div class="card score-card">
  <div class="match-header">
    <div class="player-side">
      <h3>{match.player1_name || '待定'}</h3>
    </div>
    <div class="vs">VS</div>
    <div class="player-side">
      <h3>{match.player2_name || '待定'}</h3>
    </div>
  </div>

  {#if error}
    <div class="alert alert-danger">{error}</div>
  {/if}

  <h3 style="margin:20px 0 12px;">各局比分 (三局两胜/七局四胜)</h3>
  <table>
    <thead>
      <tr>
        <th>局次</th>
        <th>{match.player1_name}</th>
        <th>{match.player2_name}</th>
      </tr>
    </thead>
    <tbody>
      {#each games as game, i}
        <tr>
          <td>第{i + 1}局</td>
          <td><input type="number" min="0" bind:value={game.score1} style="width:80px;" /></td>
          <td><input type="number" min="0" bind:value={game.score2} style="width:80px;" /></td>
        </tr>
      {/each}
    </tbody>
  </table>

  <div style="display:flex;gap:12px;margin-top:24px;justify-content:center;">
    <button class="btn btn-primary" on:click={submitScore} disabled={saving}>
      {saving ? '提交中...' : '提交比分'}
    </button>
    <button class="btn btn-danger" on:click={() => declareWalkover(match.player1_id)}>
      {match.player1_name} 弃权
    </button>
    <button class="btn btn-danger" on:click={() => declareWalkover(match.player2_id)}>
      {match.player2_name} 弃权
    </button>
  </div>
</div>
{/if}

<style>
  .score-card { max-width: 600px; margin: 0 auto; }
  .match-header { display: flex; align-items: center; justify-content: center; gap: 24px; padding: 16px; }
  .player-side { text-align: center; flex: 1; }
  .player-side h3 { font-size: 18px; }
  .vs { font-size: 20px; font-weight: 700; color: var(--text-secondary); }
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; }
  .alert-danger { background: #fce8e6; color: #c5221f; }
</style>
