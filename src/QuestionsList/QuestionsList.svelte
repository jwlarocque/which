<script>
	import QuestionsListItem from "./QuestionsListItem.svelte"

	// TODO: put this in onMount
    let questions = getQuestions();

	async function getQuestions() {
		const res = await fetch("qs");
		const data = await res.json();

		if (res.ok) {
			return data;
		} else {
			throw new Error(data);
		}
	}
</script>

<style>
	main {
		text-align: left;
	}

	h3 {
		color: #ee4035;
	}
</style>

<main>
	<h3>Your Questions:</h3>
	{#await questions}
		<p>Loading...</p>
	{:then questions}
		{#each questions as question}
			<QuestionsListItem q={question}/>
		{/each}
	{:catch error}
		<p>{error.message}</p>
	{/await}
</main>