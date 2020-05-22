<script>
    export let results = [];
    export let q;

    let maxVotes = 0;

    $: data = results.map((result, i) => ({text: q.options[result.option_id].text, votes: result.num_votes}))
    $: maxVotes = data.reduce((r, e) => (Math.max(r, e.votes)), 0);
</script>

<style>
    .barChart {
        display: grid;
        grid-row-gap: 0.4em;
        grid-column-gap: 0.4em;
        grid-template-columns: 1fr 3fr;
    }

    .barChart p {
        grid-column: 1;
        margin: auto 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .barChart div {
        box-sizing: border-box;
        grid-column: 2;
        background-color: #242020;
        color: #eef2f3;
        text-align: right;
        padding: 0.4em;
        height: 1.8em;
        margin: auto 0;
        min-width: 1.4em;
        transition: 0.6s ease-in-out;
    }
</style>

{#if data && data.length && q}
    <div class="barChart">
        {#each data as datum, i}
            <p style={"grid-row: " + (i + 1).toString() + ";"} title={datum.text}>{datum.text}</p>
            <div class="indicator" style={"grid-row: " + (i + 1).toString() + "; width: " + (100 * datum.votes / maxVotes) + "%;" + (datum.votes == maxVotes ? " background-color: #ee4035;" : "")}>
                <p>{datum.votes}</p>
            </div>
        {/each}
    </div>
{/if}
