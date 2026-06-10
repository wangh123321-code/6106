<script>
  import { api, getCurrentUser } from '$lib/api';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  let eventId = '';
  let registration = null;
  let user = null;

  onMount(() => {
    user = getCurrentUser();
    eventId = $page.params.event_id;
  });

  async function handleRegister() {
    const qualUrl = document.getElementById('qual_url').value;
    if (!qualUrl) return;
    const res = await api.post(`/events/${eventId}/register`, {
      qualification_url: qualUrl
    });
    if (res.data) registration = res.data;
  }
</script>

<svelte:head>
  <title>项目报名 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>项目报名</h1>
</div>

<div class="card" style="max-width:500px;">
  {#if registration}
    <div class="alert alert-success">报名成功，请等待审核</div>
  {:else}
    <div class="form-group">
      <label>资格证明文件URL</label>
      <input id="qual_url" type="text" placeholder="上传资格证明后粘贴链接" />
    </div>
    <button class="btn btn-primary" style="width:100%;justify-content:center;" on:click={handleRegister}>
      提交报名
    </button>
  {/if}
</div>

<style>
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; }
  .alert-success { background: #e6f4ea; color: #137333; }
</style>
