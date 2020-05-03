<script>
    export let optionCounts;
    export let q;
    
    $: sortedOptions = q.options.concat().sort((a, b) => optionCounts[b.option_id] - optionCounts[a.option_id]);
</script>

<style>
    .barChart {
        display: grid;
        grid-row-gap: 0.4em;
        grid-column-gap: 0.4em;
        grid-template-columns: 1fr 5fr;
        width: 30em;
        margin: 2em auto;
        text-align: left;
    }

    .barChart p {
        grid-column: 1;
        margin: 0.4em 0;
    }

    .barChart div {
        grid-column: 2;
        background-color: #242020;
        color: #eef2f3;
        text-align: center;
        height: 1.8em;
        margin: auto 0;
        min-width: 1em !important;
    }
</style>

{#if optionCounts.length && q}
    <div class="barChart">
        {#each sortedOptions as option, i}
            <p style={"grid-row: " + (i + 1).toString() + ";"}>{option.text}</p>
            <div class="indicator" style={"grid-row: " + (i + 1).toString() + "; width: " + (100 * optionCounts[option.option_id] / Math.max.apply(null, optionCounts)) + "%;"}>
                <p>{optionCounts[option.option_id].toString()}</p>
            </div>
        {/each}
    </div>
{/if}