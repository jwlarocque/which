<script>
	import QuestionsListItem from "./QuestionsListItem.svelte"

	const MAX_OPTIONS = 50;
	const OPTION_PLACEHOLDERS = ["Apricot", "Rhubarb", "Sour Cherry", "Raspberry", "Hot Pepper", "Gooseberry", "Peach", "Quince", "Lingonberry", "Quince", "Cloudberry", "", "Strawberry", "Blackberry", "Blueberry", "Grape", "Orange Marmalade", "Plum", "Apple Butter", "Fig"]; // example options

	// this is the kind of function Gotham deserves
	function randomJam() {
		return "e.g., " + OPTION_PLACEHOLDERS[Math.floor(Math.random() * OPTION_PLACEHOLDERS.length)]
	}

	let newQName = "";
	let newQOptions = [{id: 0, text: ""}, {id: 1, text: ""}, {id: 2, text: ""}];

	newQOptions.forEach(q => q.placeholder = randomJam());

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

	function handleNewQuestion() {
		// TODO: send request to server
		console.log("new question!");
	}

	function handleOptionUpdate() {
		if (newQOptions[newQOptions.length - 1].text !== "" && newQOptions.length < MAX_OPTIONS) {
			newQOptions = [...newQOptions, {id: newQOptions.length, text: "", placeholder: randomJam()}]
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

	form {
		display: flexbox;
		padding: 1em;
		max-width: 20em;
		margin: 0 auto;
	}

	form * {
		width: 100%;
	}

	#formBackground {
		width: 100%;
		background-color: #242020;
	}

	form input {
		background-color: #eef2f3;
		color: #242020;
		border: none;
		border-radius: 2px;
	}

	form label {
		color: #eef2f3;
		margin: 1em 0;
	}

	div {
		margin: 0.5em 0;
	}
</style>

<main>
	<div id="formBackground">
		<form id="newQForm" on:submit|preventDefault={handleNewQuestion}>
			<div class="formItem">
				<label class="formLabel"><h3>New Question</h3></label>
				<input id="newQName" placeholder="e.g., Which jam would you prefer?" bind:value={newQName}>
			</div>
			<label>Options</label>
			{#each newQOptions as option}
				<div class="newQOption">
					<input bind:value={option.text} placeholder={option.placeholder} on:input={handleOptionUpdate}>
				</div>
			{/each}
			<div>
				<button type=submit class="clickable">
					Create
				</button>
			</div>
		</form>
	</div>
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