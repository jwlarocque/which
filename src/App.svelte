<script>
	import NewQuestion from './QuestionsList/NewQuestion.svelte'
	import QuestionsList from './QuestionsList/QuestionsList.svelte'
	
	let authed = "pending"
	getAuthStatus();
	async function getAuthStatus() {
		const res = await fetch("auth/status");
		const data = await res.json();

		if (res.ok) {
			authed = data.authed;
		} else {
			throw new Error(data);
		}
	}
</script>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		background-color: #eef2f3;
		font-size: 1.2em;
		font-family: "Futura", "Ubuntu", "Helvetica Neue", "sans-serif";
	}

	:global(.clickable) {
		cursor: pointer;
		user-select: none;
	}

	:global(.visible) {
		visibility: visible;
	}
	
	:global(.hidden) {
		visibility: hidden;
	}

	:global(.button, form button) {
		border: none;
		border-radius: 2px;
		color: #eef2f3 !important;
		text-decoration: none;
		background-color: #445261;
		padding: 0.5em;
	}

	main {
		text-align: center;
		padding: 0;
		height: 100%;
		margin: 0 auto;
	}

	h1 {
		color: #ee4035;
		font-size: 4em;
		font-weight: 100;
		margin: 0;
		padding-top: 1em;
		padding-bottom: 0.5em;
	}

	p {
		max-width: 30em;
		margin: 1em auto;
		padding: 0 0.5em;
	}

	@media (max-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>

<main>
	<h1>Which?</h1>
	<p>You wouldn't use simple plurality voting to elect your government, so why would you use it to choose a jam flavor? <br/> <a href="https://en.wikipedia.org/wiki/Plurality_voting">Wait...</a></p>
	{#if authed === "true"}
		<NewQuestion/>
		<QuestionsList/>
		<a class="button" href="auth/logout">Log Out</a>
	{:else if authed === "false"}
		<br/>
		<a class="button" href="auth/login">Log In with Google</a>
	{:else}
		<p>Loading...</p>
	{/if}
</main>