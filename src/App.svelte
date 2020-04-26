<script>
	import NewQuestion from './NewQuestion/NewQuestion.svelte'
	import QuestionsList from './QuestionsList/QuestionsList.svelte'
	import {auth_state} from "./stores.js"

	let authed;

	const unsubscribe = auth_state.subscribe(value => {
		authed = value;
	})
	
	getAuthStatus();
	async function getAuthStatus() {
		const res = await fetch("auth/status");
		const data = await res.json();

		if (res.ok) {
			auth_state.set(data.authed);
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

	:global(p a, p a:visited) {
		color: #ee4035;
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
		transition: color 0.1s ease-in-out;
	}

	:global(.button:hover, form button:hover) {
		color: white !important;
	}

	:global(.radioSelect) {
		display: inline-flex;
		margin: 0;
	}

	:global(.radioSelect input) {
		opacity: 0;
		position: fixed;
		width: 0;
		z-index: -10;
		padding: 0;
		margin: 0;
		border: 0;
	}

	:global(.radioSelect label) {
		text-align: center;
		color: #eef2f3;
		background-color: #445261;
		border: 1px solid #445261;
		border-radius: 0;
		padding: 0.5em;
		margin: 0;
		user-select: none;
		display: flex;
		align-items: center;
		transition: color 0.1s ease-in-out;
	}

	:global(.radioSelect label:hover) {
		color: white;
	}

	:global(.radioSelect label > p) {
		margin: auto;
	}

	:global(.radioSelect label:first-of-type) {
		border-top-left-radius: 2px;
		border-bottom-left-radius: 2px;
	}

	:global(.radioSelect label:last-of-type) {
		border-top-right-radius: 2px;
		border-bottom-right-radius: 2px;
	}

	:global(.radioSelect input:focus + label) {
		border: 1px solid #eef2f3;
	}

	:global(.radioSelect input:checked + label) {
		background-color: #ee4035;
		border-color: #ee4035;
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
	{#if authed === "true"}
		<NewQuestion/>
		<QuestionsList/>
		<a class="button" href="auth/logout">Log Out</a>
	{:else if authed === "false"}
		<br/>
		<a class="button" href={"auth/login/" + window.location.search}>Log In with Google</a>
	{:else}
		<p>Loading...</p>
	{/if}
</main>