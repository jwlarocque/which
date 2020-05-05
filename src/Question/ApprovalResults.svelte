<script>
    export let optionCounts;
    export let q;

    let maxVotes = 0;

    $: sortedOptions = q.options.concat().sort((a, b) => optionCounts[b.option_id] - optionCounts[a.option_id]);
    $: data = sortedOptions.map((option, i) => ({text: option.text, votes: optionCounts[option.option_id]}));
    $: maxVotes = data.reduce((r, e) => (Math.max(r, e.votes)), 0);
</script>

<style>
    .barChart {
        display: grid;
        grid-row-gap: 0.4em;
        grid-column-gap: 0.4em;
        grid-template-columns: 1fr 3fr;
        max-width: 30em;
        margin: 2em auto;
        text-align: left;
        padding: 1em;
    }

    .barChart p {
        grid-column: 1;
        margin: auto 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .barChart div {
        grid-column: 2;
        background-color: #242020;
        color: #eef2f3;
        text-align: right;
        padding: 0.4em;
        height: 1em;
        margin: auto 0;
        min-width: 0.6em;
        transition: 0.6s ease-in-out;
    }
</style>

{#if optionCounts.length && q}
    <div class="barChart">
        {#each data as datum, i}
            <p style={"grid-row: " + (i + 1).toString() + ";"} title={datum.text}>{datum.text}</p>
            <div class="indicator" style={"grid-row: " + (i + 1).toString() + "; width: " + (100 * datum.votes / maxVotes) + "%;"}>
                <p>{datum.votes}</p>
            </div>
        {/each}
    </div>
{/if}
