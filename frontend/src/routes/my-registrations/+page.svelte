<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let user = null;
  let myMatches = [];
  let showAppealModal = false;
  let selectedMatch = null;
  let appealForm = { reason: '', evidence: '' };
  let submitting = false;

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    user = getCurrentUser();
    if (user?.role !== 'player') { goto('/'); return; }
  });

  function openAppeal(match) {
    selectedMatch = match;
    appealForm = { reason: '', evidence: '' };
    showAppealModal = true;
  }

  async function submitAppeal() {
    if (!appealForm.reason || !appealForm.evidence) {
      alert('请填写申诉理由和证据描述');
      return;
    }
    submitting = true;
    try {
      const res = await api.post('/appeals', {
        match_id: selectedMatch.id,
        reason: appealForm.reason,
        evidence: appealForm.evidence
      });
      if (res.code === 0) {
        alert('申诉提交成功');
        showAppealModal = false;
        selectedMatch = null;
      } else {
        alert(res.message || '提交失败');
      }
    } finally {
      submitting = false;
    }
  }
</script>

<svelte:head>
  <title>我的报名 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>我的报名</h1>
  <a href="/appeals" class="btn btn-primary">我的申诉</a>
</div>

<div class="card">
  <h3 style="margin-bottom:16px;">快速操作</h3>
  <div style="display:flex;gap:12px;flex-wrap:wrap;">
    <a href="/tournaments" class="btn btn-outline">浏览赛事</a>
    <button class="btn btn-primary" on:click={() => showAppealModal = true}>提交成绩申诉</button>
  </div>
</div>

<div class="card" style="text-align:center;padding:32px;">
  <p style="color:var(--text-secondary);">您可以通过赛事详情页进行报名</p>
  <a href="/tournaments" class="btn btn-primary" style="margin-top:12px;">浏览赛事</a>
</div>

{#if showAppealModal}
  <div class="modal-overlay" on:click={() => showAppealModal = false}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h3>提交成绩申诉</h3>
        <button class="btn btn-outline" on:click={() => showAppealModal = false}>×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label>比赛ID</label>
          <input type="text" bind:value={selectedMatch?.id || ''} placeholder="请输入比赛ID" />
        </div>
        <div class="form-group">
          <label>申诉理由</label>
          <textarea rows="3" bind:value={appealForm.reason} placeholder="请详细说明申诉理由"></textarea>
        </div>
        <div class="form-group">
          <label>证据描述</label>
          <textarea rows="3" bind:value={appealForm.evidence} placeholder="请提供相关证据描述"></textarea>
        </div>
        <p style="font-size:12px;color:var(--text-secondary);">
          注意：需在成绩发布后48小时内提交申诉，且您必须为该场比赛的参赛选手。
        </p>
      </div>
      <div class="modal-footer">
        <button class="btn btn-outline" on:click={() => showAppealModal = false}>取消</button>
        <button class="btn btn-primary" on:click={submitAppeal} disabled={submitting}>
          {submitting ? '提交中...' : '提交申诉'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed; inset: 0; background: rgba(0,0,0,0.5);
    display: flex; align-items: center; justify-content: center; z-index: 1000;
  }
  .modal {
    background: var(--surface); border-radius: var(--radius);
    width: 90%; max-width: 560px;
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
