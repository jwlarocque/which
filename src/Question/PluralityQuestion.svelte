<script>
    import {scale} from 'svelte/transition';
    export let q;
    export let votes;

    let selectedOption = -1;
    $: if (selectedOption == -1) {selectedOption = votes.indexOf(1);} else if (selectedOption >= 0) {votes = votes.fill(0); votes[selectedOption] = 1;}
</script>

<script context="module">
    export function pluralityVotesFromBallot(q, ballot) {
        var vs = new Array(q.options.length);
        if (ballot && ballot.votes && ballot.votes.length) {
            ballot.votes.forEach(vote => (vs[vote.option_id] = vote.state));
        }
        return vs
    }
</script>

<style>
    label {
        position: relative;
        display: flex;
        align-items: center;
    }

    input[type="radio"] {
        position: absolute;
        opacity: 0;
    }

    input[type="radio"] + span {
        box-sizing: border-box;
        border: 1px solid #eef2f3;
        border-radius: 3.2em;

        /* this is what it takes to make css respect you */
        height: 1.6em;
        max-height: 1.6em;
        min-height: 1.6em;
        width: 1.6em;
        max-width: 1.6em;
        min-width: 1.6em;

        margin: 0.6em;
        display: inline-flex;
        user-select: none;

        transition: 0.1s ease-in-out;
    }

    input[type="radio"]:checked + span {
        border-width: 8px;
    }
</style>

{#if q.name}
    <!-- TODO: better explainer -->
    <p>Select all options you are okay with.</p>
    {#each q.options as option}
        <label class="clickable">
            <input type="radio" bind:group={selectedOption} value={option.option_id}>
            <span></span>
            {option.text}
        </label>
    {/each}
{:else}
    <p>Loading...</p>
{/if}