<script>
  import { api, isLoggedIn } from '$lib/api';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let registrations = [];
  let selectedIds = new Set();

  $: eventId = $page.params.event_id;

  onMount(async () => {
    await loadRegistrations();
  });

  async function loadRegistrations() {
    const res = await api.get(`/events/${eventId}/registrations`);
    if (res.data) registrations = res.data;
  }

  function toggleSelect(id) {
    if (selectedIds.has(id)) {
      selectedIds.delete(id);
    } else {
      selectedIds.add(id);
    }
    selectedIds = selectedIds;
  }

  async function batchApprove() {
    if (selectedIds.size === 0) return;
    const res = await api.post(`/events/${eventId}/batch-approve`, {
      registration_ids: Array.from(selectedIds)
    });
    if (res.code === 0) {
      selectedIds = new Set();
      await loadRegistrations();
    }
  }

  async function rejectReg(id) {
    await api.put(`/registrations/${id}/reject`);
    await loadRegistrations();
  }

  function statusLabel(s) {
    const map = { pending: '待审核', approved: '已通过', rejected: '已拒绝' };
    return map[s] || s;
  }

  function statusClass(s) {
    const map = { pending: 'badge-warning', approved: 'badge-success', rejected: 'badge-danger' };
    return map[s] || '';
  }
</script>

<svelte:head>
  <title>报名审核 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>报名审核</h1>
  <div style="display:flex;gap:8px;">
    <button class="btn btn-success" on:click={batchApprove} disabled={selectedIds.size === 0}>
      批量通过 ({selectedIds.size})
    </button>
    <button class="btn btn-outline" on:click={() => api.post(`/events/${eventId}/close-registration`).then(loadRegistrations)}>
      关闭报名
    </button>
    <button class="btn btn-primary" on:click={() => api.post(`/events/${eventId}/draw`).then(() => {})}>
      一键抽签
    </button>
  </div>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th><input type="checkbox" on:change={(e) => {
          if (e.target.checked) {
            registrations.filter(r => r.status === 'pending').forEach(r => selectedIds.add(r.id));
          } else {
            selectedIds = new Set();
          }
        }} /></th>
        <th>选手</th>
        <th>资格证明</th>
        <th>报名时间</th>
        <th>状态</th>
        <th>操作</th>
      </tr>
    </thead>
    <tbody>
      {#each registrations as r}
        <tr>
          <td>
            {#if r.status === 'pending'}
              <input type="checkbox" checked={selectedIds.has(r.id)} on:change={() => toggleSelect(r.id)} />
            {/if}
          </td>
          <td>{r.player_id}</td>
          <td><a href={r.qualification_url} target="_blank">查看</a></td>
          <td>{r.created_at?.substring(0,10)}</td>
          <td><span class="badge {statusClass(r.status)}">{statusLabel(r.status)}</span></td>
          <td>
            {#if r.status === 'pending'}
              <button class="btn btn-danger" style="padding:2px 10px;font-size:12px;" on:click={() => rejectReg(r.id)}>拒绝</button>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
