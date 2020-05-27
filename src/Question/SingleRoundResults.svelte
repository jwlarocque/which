<script>
    import { flip } from 'svelte/animate';

    export let results = [];
    export let q;

    let maxVotes = 0;

    $: data = results.map((result, i) => ({id: result.option_id, text: q.options[result.option_id].text, votes: result.num_votes}))
    $: maxVotes = data.reduce((r, e) => (Math.max(r, e.votes)), 0);
    $: data = data.sort((a, b) => b.votes - a.votes);
</script>

<style>
    .barChart {
        display: flex;
        flex-direction: column;
    }

    .barChart div {
        display: inline-block;
        margin-bottom: 0.6em;
    }

    .barChart p {
        margin: auto 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .barChart .indicator {
        box-sizing: border-box;
        background-color: #242020;
        color: #eef2f3;
        text-align: right;
        padding: 0.4em;
        height: 1.8em;
        margin: auto 0;
        min-width: 1.4em;
        transition: 0.4s ease-in-out;
    }
</style>

{#if data && data.length && q}
    <div class="barChart">
        {#each data as datum, i (datum.id)}
            <div animate:flip="{{duration: 400}}">
                <p title={datum.text}>{datum.text}</p>
                <div class="indicator" style={"grid-row: " + (i + 1).toString() + "; width: " + (100 * datum.votes / maxVotes) + "%;" + (datum.votes == maxVotes ? " background-color: #ee4035;" : "")}>
                    <p>{datum.votes}</p>
                </div>
            </div>
        {/each}
    </div>
{/if}
