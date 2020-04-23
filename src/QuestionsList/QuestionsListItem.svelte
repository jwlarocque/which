<script>
    // TODO: consider just putting this in QuestionsList.svelte
    export let q;
</script>

<style>
    div {
        padding: 0.2em 0;
    }

</style>

<div>
    <p>The question is: {q.name}</p>
</div>
