<script>
  import { api, getCurrentUser } from '$lib/api';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let qualificationURL = '';
  let partnerID = '';
  let eventId = '';
  let eventType = '';
  let error = '';
  let success = false;

  $: id = $page.params.id;

  onMount(() => {
    eventId = $page.params.event_id || id;
  });

  async function handleRegister() {
    error = '';
    success = false;
    const res = await api.post(`/events/${eventId}/register`, {
      qualification_url: qualificationURL,
      partner_id: partnerID || undefined
    });
    if (res.code === 0) {
      success = true;
    } else {
      error = res.message || '报名失败';
    }
  }
</script>

<svelte:head>
  <title>报名参赛 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>报名参赛</h1>
</div>

<div class="card" style="max-width:500px;">
  {#if success}
    <div class="alert alert-success">报名成功，请等待组委会审核</div>
  {/if}
  {#if error}
    <div class="alert alert-danger">{error}</div>
  {/if}

  <div class="form-group">
    <label>资格证明文件URL</label>
    <input type="text" bind:value={qualificationURL} placeholder="请上传资格证明后粘贴链接" required />
  </div>

  <div class="form-group">
    <label>搭档ID（仅混双项目）</label>
    <input type="text" bind:value={partnerID} placeholder="混合双打时填写搭档用户ID" />
  </div>

  <button class="btn btn-primary" style="width:100%;justify-content:center;" on:click={handleRegister}>
    提交报名
  </button>
</div>

<style>
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; }
  .alert-danger { background: #fce8e6; color: #c5221f; }
  .alert-success { background: #e6f4ea; color: #137333; }
</style>
