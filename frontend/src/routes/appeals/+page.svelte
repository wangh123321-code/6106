<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let user = null;
  let appeals = [];
  let filterStatus = '';
  let selectedAppeal = null;
  let reviewDecision = 'rejected';
  let reviewNote = '';
  let games = [{ score1: 0, score2: 0 }];
  let submitting = false;
  let showCreateModal = false;
  let newAppeal = { match_id: '', reason: '', evidence: '' };

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    user = getCurrentUser();
    await loadAppeals();
  });

  async function loadAppeals() {
    let path = '/appeals?';
    if (filterStatus) path += `status=${filterStatus}`;
    const res = await api.get(path);
    if (res.data) appeals = res.data;
  }

  function getStatusBadge(status) {
    const map = {
      pending: { cls: 'badge-warning', text: '待处理' },
      reviewing: { cls: 'badge-info', text: '审核中' },
      upheld: { cls: 'badge-success', text: '改判' },
      rejected: { cls: 'badge-danger', text: '维持原判' }
    };
    return map[status] || { cls: 'badge-info', text: status };
  }

  function openReview(appeal) {
    selectedAppeal = appeal;
    reviewDecision = 'rejected';
    reviewNote = '';
    games = [{ score1: 0, score2: 0 }];
  }

  function addGame() {
    games = [...games, { score1: 0, score2: 0 }];
  }

  function removeGame(i) {
    games = games.filter((_, idx) => idx !== i);
  }

  async function submitReview() {
    if (!selectedAppeal) return;
    submitting = true;
    try {
      const body = {
        decision: reviewDecision,
        review_note: reviewNote
      };
      if (reviewDecision === 'upheld') {
        body.games = games;
      }
      const res = await api.put(`/appeals/${selectedAppeal.id}/review`, body);
      if (res.code === 0) {
        alert('审核成功');
        selectedAppeal = null;
        await loadAppeals();
      } else {
        alert(res.message || '审核失败');
      }
    } finally {
      submitting = false;
    }
  }

  async function submitAppeal() {
    if (!newAppeal.match_id || !newAppeal.reason || !newAppeal.evidence) {
      alert('请填写完整申诉信息');
      return;
    }
    submitting = true;
    try {
      const res = await api.post('/appeals', newAppeal);
      if (res.code === 0) {
        alert('申诉提交成功');
        showCreateModal = false;
        newAppeal = { match_id: '', reason: '', evidence: '' };
        await loadAppeals();
      } else {
        alert(res.message || '提交失败');
      }
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>成绩申诉 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>成绩申诉</h1>
  <div style="display:flex;gap:8px;">
    {#if user?.role === 'player'}
      <button class="btn btn-primary" on:click={() => showCreateModal = true}>提交申诉</button>
    {/if}
    <button class="btn btn-outline" on:click={loadAppeals}>刷新</button>
  </div>
</div>

<div class="card" style="margin-bottom:16px;">
  <div style="display:flex;gap:12px;align-items:center;">
    <div class="form-group" style="margin:0;">
      <label>状态筛选</label>
      <select bind:value={filterStatus} onchange={loadAppeals}>
        <option value="">全部</option>
        <option value="pending">待处理</option>
        <option value="reviewing">审核中</option>
        <option value="upheld">改判</option>
        <option value="rejected">维持原判</option>
      </select>
    </div>
  </div>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>申诉编号</th>
        <th>比赛ID</th>
        <th>申诉人</th>
        <th>申诉理由</th>
        <th>状态</th>
        <th>提交时间</th>
        {#if user?.role === 'committee'}
          <th>操作</th>
        {/if}
      </tr>
    </thead>
    <tbody>
      {#each appeals as appeal}
        <tr>
          <td style="font-family:monospace;font-size:12px;">{appeal.appeal_id}</td>
          <td style="font-family:monospace;font-size:12px;">{appeal.match_id}</td>
          <td>{appeal.appellant_id}</td>
          <td style="max-width:240px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" title={appeal.reason}>{appeal.reason}</td>
          <td>
            <span class={`badge ${getStatusBadge(appeal.status).cls}`}>
              {getStatusBadge(appeal.status).text}
            </span>
          </td>
          <td>{appeal.created_at?.substring(0,19).replace('T', ' ')}</td>
          {#if user?.role === 'committee'}
            <td>
              {#if appeal.status === 'pending' || appeal.status === 'reviewing'}
                <button class="btn btn-primary btn-sm" on:click={() => openReview(appeal)}>审核</button>
              {/if}
            </td>
          {/if}
        </tr>
      {/each}
      {#if appeals.length === 0}
        <tr>
          <td colspan={user?.role === 'committee' ? 7 : 6} style="text-align:center;padding:32px;color:var(--text-secondary);">
            暂无申诉记录
          </td>
        </tr>
      {/if}
    </tbody>
  </table>
</div>

{#if selectedAppeal}
  <div class="modal-overlay" on:click={() => selectedAppeal = null}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h3>审核申诉 - {selectedAppeal.appeal_id}</h3>
        <button class="btn btn-outline" on:click={() => selectedAppeal = null}>×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>比赛ID</label>
          <input type="text" value={selectedAppeal.match_id} disabled />
        </div>
        <div class="form-group">
          <label>申诉理由</label>
          <textarea rows="2" disabled>{selectedAppeal.reason}</textarea>
        </div>
        <div class="form-group">
          <label>证据描述</label>
          <textarea rows="3" disabled>{selectedAppeal.evidence}</textarea>
        </div>
        <div class="form-group">
          <label>审核决定</label>
          <select bind:value={reviewDecision}>
            <option value="rejected">维持原判</option>
            <option value="upheld">改判</option>
          </select>
        </div>
        {#if reviewDecision === 'upheld'}
          <div class="form-group">
            <label>新比分（每局得分）</label>
            {#each games as g, i}
              <div style="display:flex;gap:8px;margin-bottom:8px;align-items:center;">
                <span>第{i + 1}局:</span>
                <input type="number" bind:value={g.score1} placeholder="选手1得分" style="width:120px;" />
                <span>VS</span>
                <input type="number" bind:value={g.score2} placeholder="选手2得分" style="width:120px;" />
                {#if games.length > 1}
                  <button class="btn btn-outline btn-sm" on:click={() => removeGame(i)}>删除</button>
                {/if}
              </div>
            {/each}
            <button class="btn btn-outline btn-sm" on:click={addGame}>+ 添加局数</button>
          </div>
        {/if}
        <div class="form-group">
          <label>审核备注</label>
          <textarea rows="2" bind:value={reviewNote} placeholder="请填写审核备注"></textarea>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-outline" on:click={() => selectedAppeal = null}>取消</button>
        <button class="btn btn-primary" on:click={submitReview} disabled={submitting}>
          {submitting ? '提交中...' : '确认审核'}
        </button>
      </div>
    </div>
  </div>
{/if}

{#if showCreateModal}
  <div class="modal-overlay" on:click={() => showCreateModal = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h3>提交成绩申诉</h3>
        <button class="btn btn-outline" on:click={() => showCreateModal = false}>×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>比赛ID</label>
          <input type="text" bind:value={newAppeal.match_id} placeholder="请输入比赛ID" />
        </div>
        <div class="form-group">
          <label>申诉理由</label>
          <textarea rows="3" bind:value={newAppeal.reason} placeholder="请详细说明申诉理由"></textarea>
        </div>
        <div class="form-group">
          <label>证据描述</label>
          <textarea rows="3" bind:value={newAppeal.evidence} placeholder="请提供相关证据描述"></textarea>
        </div>
        <p style="font-size:12px;color:var(--text-secondary);">
          注意：需在成绩发布后48小时内提交申诉，且您必须为该场比赛的参赛选手。
        </p>
      </div>
      <div class="modal-footer">
        <button class="btn btn-outline" on:click={() => showCreateModal = false}>取消</button>
        <button class="btn btn-primary" on:click={submitAppeal} disabled={submitting}>
          {submitting ? '提交中...' : '提交申诉'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .btn-sm { padding: 4px 12px; font-size: 12px; }
  .modal-overlay {
    position: fixed; inset: 0; background: rgba(0,0,0,0.5);
    display: flex; align-items: center; justify-content: center; z-index: 1000;
  }
  .modal {
    background: var(--surface); border-radius: var(--radius);
    width: 90%; max-width: 560px; max-height: 90vh; overflow-y: auto;
  }
  .modal-header {
    display: flex; justify-content: space-between; align-items: center;
    padding: 16px 20px; border-bottom: 1px solid var(--border);
  }
  .modal-header h3 { margin: 0; font-size: 18px; }
  .modal-body { padding: 20px; }
  .modal-footer {
    padding: 16px 20px; border-top: 1px solid var(--border);
    display: flex; justify-content: flex-end; gap: 8px;
  }
</style>
