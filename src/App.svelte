<script>
	import {onMount} from 'svelte';
	import NewQuestion from './NewQuestion/NewQuestion.svelte';
	import QuestionsList from './QuestionsList/QuestionsList.svelte';
	import Question from './Question/Question.svelte';
	import {auth_state} from "./stores.js";

	let authed;
	let question_id = "";

	const unsubscribe = auth_state.subscribe(value => {
		authed = value;
	})
	
	getAuthStatus();
	async function getAuthStatus() {
		const res = await fetch("auth/status");
		const data = await res.json();

		if (res.ok) {
			auth_state.set(data.authed);
			if (authed === "true") {
				checkQuery();
			}
		} else {
			throw new Error(data);
		}
	}

	function checkQuery() {
		let params = (new URL(document.location)).searchParams;
		question_id = params.get("q") || "";
	}
</script>

<style>
	/* TODO: better global css (use preprocessor?) */
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

	:global(.darkBackground) {
		width: 100%;
        background-color: #242020;
		color: #eef2f3;
        position: relative;
	}

	main {
		text-align: center;
		padding: 0;
		height: 100%;
		margin: 0 auto;
	}

	@media (max-width: 640px) {
		main {
			max-width: none;
		}
	}

	h1 {
		color: #ee4035;
		font-size: 4em;
		font-weight: 100;
		margin: 0;
		padding: 0.5em;
	}

	h1 a {
		text-decoration: inherit;
		color: inherit;
	}

	p {
		max-width: 30em;
		margin: 1em auto;
		padding: 0 0.5em;
	}

	nav {
		width: calc(100% - 2em);
		max-width: 30em;
		display: inline-flex;
		justify-content: space-between;
		flex-wrap: wrap;
		row-gap: 1em;
		margin: 2em 1em;
		box-sizing: border-box;
	}

	#navLinks {
		display: inline-flex;
	}

	nav * {
		margin: auto 0;
	}

	nav a, nav a:visited {
		color: #242020;
	}

	nav .button {
		box-sizing: border-box;
	}

</style>

<svelte:head>
  <title>Which?</title>
  <link rel="icon" type="image/png" href="favicon.ico"/>
</svelte:head>

<main>
	<nav>
		<div id="navLinks">
			<a href="http://jwlarocque.com/">jwlarocque.com</a>
			<p>/</p>
			<a href="https://github.com/jwlarocque/which">Source</a>
		</div>
		<div>
			{#if authed === "true"}
				<a class="button" href="auth/logout">Log Out</a>
			{:else if authed === "false"}
				<a class="button" href={"auth/login/" + window.location.search}>Log In</a>
			{:else}
				<p>Loading...</p>
			{/if}
		</div>
	</nav>
	<h1><a href="/">Which?</a></h1>
	{#if authed === "true"}
		{#if question_id.length > 0}
			<Question id={question_id}/>
		{:else}
			<NewQuestion/>
			<QuestionsList/>
		{/if}
	{:else if authed === "false"}
		<br/>
		<a class="button" href={"auth/login/" + window.location.search}>Log In with Google</a>
	{:else}
		<p>Loading...</p>
	{/if}
</main>