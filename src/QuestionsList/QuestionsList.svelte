<script>
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

	function handleClick() {
		questions = getQuestions();
	}
</script>

<main>
    <button on:click={handleClick}>
		Get questions
	</button>
	{#await questions}
		<p>Loading...</p>
	{:then questions}
		{#each questions as question}
			<p>{question.Name}</p>
		{/each}
	{:catch error}
		<p>{error.message}</p>
	{/await}
</main>