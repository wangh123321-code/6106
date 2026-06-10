<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let matches = [];
  let user = null;

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    user = getCurrentUser();
    if (user?.role !== 'referee') { goto('/'); return; }
    await loadMatches();
  });

  async function loadMatches() {
    const res = await api.get('/referee/matches');
    if (res.data) matches = res.data;
  }

  function statusLabel(s) {
    const map = { pending: '待开始', completed: '已结束', walkover: '弃权' };
    return map[s] || s;
  }

  function statusClass(s) {
    const map = { pending: 'badge-warning', completed: 'badge-success', walkover: 'badge-danger' };
    return map[s] || '';
  }
</script>

<svelte:head>
  <title>裁判工作台 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>裁判工作台</h1>
  <button class="btn btn-outline" on:click={loadMatches}>刷新</button>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>轮次</th>
        <th>选手1</th>
        <th>比分</th>
        <th>选手2</th>
        <th>状态</th>
        <th>操作</th>
      </tr>
    </thead>
    <tbody>
      {#each matches as m}
        <tr>
          <td>第{m.round}轮</td>
          <td>{m.player1_name || '待定'}</td>
          <td>{m.status === 'pending' ? 'vs' : `${m.score1} : ${m.score2}`}</td>
          <td>{m.player2_name || '待定'}</td>
          <td><span class="badge {statusClass(m.status)}">{statusLabel(m.status)}</span></td>
          <td>
            {#if m.status === 'pending' && m.player1_id && m.player2_id}
              <a href="/matches/{m.id}/score" class="btn btn-primary" style="padding:4px 12px;">录入比分</a>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
