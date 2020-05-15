<script>
	import {fade} from 'svelte/transition';
	import {qs} from "../stores.js";

	// TODO: put this in onMount
	let questions;
	const unsubscribe = qs.subscribe(value => {
		questions = value;
	});

	getQuestions();

	async function getQuestions() {
		const res = await fetch("qs/list");
		const data = await res.json();

		if (res.ok) {
			qs.set(data);
		} else {
			throw new Error(data);
		}
	}

	function fadeOutIn(element) {
		element.animate([
			{opacity: 1.0},
			{opacity: 0.0},
			{opacity: 0.0},
			{opacity: 1.0},
			{easing: ["ease-in-out"]}],
			1000);
	}

	function copyToClipboard(link, element) {
		navigator.clipboard.writeText(link).then(fadeOutIn(element));
	}

	async function deleteQuestion(q, element) {
		if (!confirm("Are you sure you want to delete the question \"" + q.name + "\" and all associated data?  This cannot be undone.")) {
			return;
		}
		console.log("deleting...");
		const res = await fetch("qs/del", {
			method: "POST",
			headers: {"Content-Type": "text/plain"},
			body: q.question_id,
		});
		
		if (res.ok) {
			qs.update((value) => value.filter((e) => e !== q));
		} else {
			alert("failed to delete question: " + res.statusText);
		}
	}
</script>

<style>
	main {
		text-align: left;
		max-width: 30em;
		margin: 0 auto;
		padding: 1em;
	}

	h3 {
		color: #ee4035;
	}

	hr {
		color: rgba(0, 0, 0, 0.2);
	}

	.questionRow {
		display: grid;
		grid-template-columns: 3fr 1fr 1fr;
		grid-column-gap: 0.4em;
	}

	.questionRow hr, .across {
		grid-column: 1/6;
		width: 100%;
	}

	.questionRow .icon {
		margin: auto;
		position: relative;
	}

	.questionRow img {
		margin: auto;
		cursor: pointer;
		background-color: #eef2f3;
		position: absolute;
		transform: translate(-50%, -50%);
	}

	.center {
		text-align: center;
		margin: auto 0;
	}

	.question * {
        text-decoration: inherit;
        color: inherit;
    }

	#noQuestions {
		margin: 2em;
		grid-column: 1/6;
		opacity: 0.6;
	}
</style>

<main>
	{#await questions}
		<p>Loading...</p>
	{:then questions}
		<div id="questionsList">
			<div class="questionRow">
				<h3>Your Questions</h3>
				<h4 class="center">Link</h4>
				<h4 class="center">Delete</h4>
			</div>
			{#if questions.length > 0}
				{#each questions as q, i}
					{#if i > 0}<hr transition:fade/>{:else}<div class="across"></div>{/if}
					<div class="questionRow" transition:fade>
						<p class="question"><a href={"/?q=" + q.question_id}>{q.name}</a></p>
						<div class="icon">
							<img class="center" src="images/done.svg" alt="copied"/>
							<img class="center" src="images/copy.svg" alt="copy to clipboard" title="copy to clipboard" on:click={copyToClipboard(window.location.host + "/?q=" + q.question_id, this)}/>
						</div>
						<div class="icon">
							<img class="center" src="images/done.svg" alt="deleted"/>
							<img class="center" src="images/delete.svg" alt="delete" title="delete" on:click={deleteQuestion(q, this)}/>
						</div>
					</div>
				{/each}
			{:else}
				<hr/>
				<p class="center" id="noQuestions">Create a new question with the form above.</p>
			{/if}
		</div>
	{:catch error}
		<p>{error.message}</p>
	{/await}
</main>