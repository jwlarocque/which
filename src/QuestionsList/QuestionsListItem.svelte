<script>
    // TODO: consider just putting this in QuestionsList.svelte
    export let q;
</script>

<style>
    div {
        padding: 0 0;
    }

    .question * {
        text-decoration: inherit;
        color: inherit;
    }
</style>

<div>
    <p class="question"><a href={"/?q=" + q.question_id}>{q.name}</a></p>
</div>
