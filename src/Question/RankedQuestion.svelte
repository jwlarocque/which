<script>
    import {onMount} from 'svelte';
    import {scale} from 'svelte/transition';
    export let q;
    export let votes;

    $: q.options.sort((a, b) => (votes[a.option_id] - votes[b.option_id]))

    function upRank(option) {
        if (votes[option.option_id] > 1) {
            votes[option.option_id]--;
        }
    }

    function downRank(option) {
        if (votes[option.option_id] < votes.length - 1) {
            votes[option.option_id]++;
        }
    }

</script>

<script context="module">
    export function rankedVotesFromBallot(q, ballot) {
        var vs = new Array(q.options.length).fill().map((_, i) => i + 1);
        if (ballot && ballot.votes) {
            ballot.votes.forEach(vote => (vs[vote.option_id] = vote.state));
        }
        return vs
    }
</script>

<style>
    .dragList {
        position: relative;    
    }

    .dragList > div {
        position: relative;
        display: flex;
        align-items: center;
    }

    label > * {
        margin-right: 1em;
    }

    .grabbed {
        /*opacity: 0.0;*/
        background-color: blue !important;
    }

    #ghost {
        position: absolute;
        background-color: red;
        left: -50px;
    }
</style>

{#if q.name}
    <!-- TODO: better explainer -->
    <p>Rank the options from most preferred to least preferred.</p>
    <div class="dragList">
        <div id="ghost"></div>
        {#each q.options as option, i}
            <div 
            id={"item" + option.option_id}>
                <!--{#if i != 0} <button on:click={upRank(option)}>^</button> {/if}-->
                <p>{votes[option.option_id]}</p>
                <!--{#if i != q.options.length - 1} <button on:click={downRank(option)}>v</button> {/if}-->
                <p>{option.text}</p>
            </div>
        {/each}
    </div>
{:else}
    <p>Loading...</p>
{/if}