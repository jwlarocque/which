<script>
    import {scale} from 'svelte/transition';
    export let q;
    export let votes;
</script>

<style>
    label {
        position: relative;
        display: flex;
        align-items: center;
    }

    input[type="checkbox"] {
        position: absolute;
        opacity: 0;
    }

    input[type="checkbox"] + span {
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

    input[type="checkbox"]:checked + span {
        background-color: #eef2f3;
    }

    input[type="checkbox"] + span > img {
        height: 1.4em;
        margin: auto;
        width: auto;
        opacity: 0.0;
        transition: 0.1s ease-in-out;
    }

    input[type="checkbox"]:checked + span > img {
        opacity: 1.0;
    }
</style>

{#if q.name}
    <!-- TODO: better explainer -->
    <p>Select all options you are okay with.</p>
    {#each q.options as option}
        <label class="clickable">
            <input type="checkbox" bind:checked={votes[option.option_id]}>
            <span><img src="images/done.svg" alt="checkmark"/></span>
            {option.text}
        </label>
    {/each}
{:else}
    <p>Loading...</p>
{/if}