<script>
    import {onMount} from 'svelte';
    import {scale} from 'svelte/transition';

    import DragDropList from 'svelte-dragdroplist';

    export let q;
    export let votes;

    $: votes = votesFromOpts(q.options);

    function votesFromOpts(opts) {
        let vs = new Array();
        opts.forEach((opt, i) => (vs[opt.option_id] = opts.length - i));
        return vs;
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

    export function orderOptsByVote(opts, votes) {
        opts.sort((a, b) => votes[b.option_id] - votes[a.option_id]);
    }
</script>

<style>
    :global(.dragdroplist .item) {
        color: #242020;
    }

    label > * {
        margin-right: 1em;
    }
</style>

{#if q.sorted}
    <!-- TODO: better explainer -->
    <p>Rank the options from most preferred to least preferred.</p>
    <DragDropList bind:data={q.options}/>
{:else}
    <p>Loading...</p>
{/if}