<script>
  export let matches = [];
  export let event = null;

  $: rounds = buildRounds(matches);
  $: numRounds = rounds.length;

  function buildRounds(matches) {
    const map = {};
    matches.forEach(m => {
      if (!map[m.round]) map[m.round] = [];
      map[m.round].push(m);
    });
    const keys = Object.keys(map).map(Number).sort((a, b) => a - b);
    return keys.map(k => {
      const round = map[k].sort((a, b) => a.position - b.position);
      return round;
    });
  }

  function roundLabel(idx) {
    if (idx === numRounds - 1) return '决赛';
    if (idx === numRounds - 2) return '半决赛';
    return `第${idx + 1}轮`;
  }

  function matchResult(m) {
    if (m.status === 'pending') return '待定';
    if (m.walkover) return '弃权晋级';
    return `${m.score1} : ${m.score2}`;
  }

  function winnerClass(m, playerNum) {
    if (!m.winner_id) return '';
    if (playerNum === 1 && m.winner_id === m.player1_id) return 'winner';
    if (playerNum === 2 && m.winner_id === m.player2_id) return 'winner';
    return 'loser';
  }
</script>

<div class="bracket-container">
  <div class="bracket-scroll">
    <div class="bracket" style="--num-rounds: {numRounds};">
      {#each rounds as round, ri}
        <div class="round" style="--round-idx: {ri};">
          <div class="round-label">{roundLabel(ri)}</div>
          <div class="round-matches" style="--match-count: {round.length};">
            {#each round as match, mi}
              <div class="match-card" style="--match-idx: {mi};">
                <div class="player {winnerClass(match, 1)}">
                  <span class="name">{match.player1_name || '待定'}</span>
                  {#if match.status !== 'pending'}
                    <span class="score">{match.score1}</span>
                  {/if}
                </div>
                <div class="player {winnerClass(match, 2)}">
                  <span class="name">{match.player2_name || '待定'}</span>
                  {#if match.status !== 'pending'}
                    <span class="score">{match.score2}</span>
                  {/if}
                </div>
                {#if match.walkover}
                  <div class="walkover-badge">弃权</div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/each}

      {#each rounds as round, ri}
        {#if ri < rounds.length - 1}
          <div class="connectors" style="--round-idx: {ri};">
            {#each round as match, mi}
              <div class="connector" style="--match-idx: {mi};"></div>
            {/each}
          </div>
        {/if}
      {/each}
    </div>
  </div>
</div>

<style>
  .bracket-container { width: 100%; overflow: hidden; }
  .bracket-scroll { overflow-x: auto; padding-bottom: 16px; }

  .bracket {
    display: flex;
    align-items: flex-start;
    gap: 0;
    min-width: max-content;
  }

  .round {
    display: flex;
    flex-direction: column;
    min-width: 200px;
  }

  .round-label {
    text-align: center;
    font-weight: 600;
    color: var(--text-secondary);
    font-size: 13px;
    padding: 8px;
    margin-bottom: 8px;
  }

  .round-matches {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    flex: 1;
    gap: 8px;
  }

  .match-card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
    position: relative;
  }

  .player {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    font-size: 13px;
    border-bottom: 1px solid var(--border);
    transition: background 0.2s;
  }

  .player:last-child { border-bottom: none; }

  .player .name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 140px; }
  .player .score { font-weight: 700; min-width: 24px; text-align: right; }

  .player.winner { background: #e6f4ea; }
  .player.winner .name { font-weight: 600; }
  .player.loser { color: var(--text-secondary); }

  .walkover-badge {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--danger);
    color: #fff;
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 10px;
    font-weight: 500;
  }

  .connectors {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    align-items: center;
    min-width: 40px;
    flex: 0 0 40px;
  }

  .connector {
    width: 100%;
    height: 2px;
    background: var(--border);
    position: relative;
  }
</style>
