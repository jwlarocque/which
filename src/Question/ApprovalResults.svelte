<script>
    export let optionCounts;
    export let q;
</script>

{#if optionCounts.length && q}
    {#each q.options as option}
        <p>{option.text}: {optionCounts[option.option_id].toString()}</p>
    {/each}
{/if}