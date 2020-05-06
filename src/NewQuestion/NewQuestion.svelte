<script>
    import {qs} from "../stores.js"

    const MAX_OPTIONS = 50;
	const OPTION_PLACEHOLDERS = ["Apricot", "Rhubarb", "Sour Cherry", "Raspberry", "Hot Pepper", "Gooseberry", "Peach", "Quince", "Lingonberry", "Quince", "Cloudberry", "Strawberry", "Blackberry", "Blueberry", "Grape", "Orange Marmalade", "Plum", "Apple Butter", "Fig"]; // example options
	const MAX_FIELD_CHARS = "128";

	// this is my kind of function
	function randomJam() {
		return "e.g., " + OPTION_PLACEHOLDERS[Math.floor(Math.random() * OPTION_PLACEHOLDERS.length)];
	}

	const newQTypes = [
		{value: "approval", text: "Approval", explainer: "<p>Good for making an everyday choice which most people are okay with.</p><p>Respondents select all the options they would be okay with, and the winner is the option the most people approve of.</p><p>More information: <a href=\"https://www.electionscience.org/library/approval-voting/\">ElectionScience.org</a></p>"}, 
		{value: "runoff", text: "Ranked Choice", explainer: "<p>Good for selecting an option with true majority support, and allows people to express their preference for choices which are unlikely to win.</p><p>Respondents rank the choices in order of preference.</p><p>More information: <a href=\"https://www.fairvote.org/rcv#where_is_ranked_choice_voting_used\">FairVote.org</a></p>"}, 
		{value: "plurality", text: "Plurality", explainer: "<p>It's straightforward I guess, if that's what you're looking for in a voting system.</p><p>Respondents pick one option; the option with the most votes wins.</p><p>More information (first section): <a href=\"https://www.fairvote.org/plurality_majority_systems\">FairVote.org</a></p>"}];
	let newQType = newQTypes[0];

    let newQFormVisible = true;
	let newQName = "";
	let newQOptions = [{option_id: 0, text: ""}, {option_id: 1, text: ""}, {option_id: 2, text: ""}];

    newQOptions.forEach(q => q.placeholder = randomJam());
    
    // TODO: better user feedback than alerts
    function validateQuestion() {
        if (newQName == "") {
            alert("New Question must have name/title.");
            return false;
        }
        if (newQOptions.filter(option => option.text.length > 0).length < 2) {
            alert("New Question must have at least two answers.");
            return false;
        }
        return true;
    }

    function clearNewQForm() {
        newQName = "";
        newQOptions.forEach(option => option.text = "");
		newQOptions = newQOptions;
    }

    async function handleNewQuestion() {
        newQFormVisible = false;
        if (validateQuestion()) {
            const res = await fetch("qs/new", {
                method: "POST",
                headers: {"Content-Type": "application/json",},
                body: JSON.stringify({"name": newQName, "type": newQType.value, "options": newQOptions}),
            });
            const data = await res.json();

            if (res.ok) {
                clearNewQForm();
                newQFormVisible = true;
                qs.update(value => [...value, data]);
            } else {
                newQFormVisible = true;
                throw new Error(data.message); // TODO: improve and replicate this error pattern
            }
        } else {
            newQFormVisible = true;
        }
	}

	function handleOptionUpdate() {
		if (newQOptions[newQOptions.length - 1].text !== "" && newQOptions.length < MAX_OPTIONS) {
			newQOptions = [...newQOptions, {option_id: newQOptions.length, text: "", placeholder: randomJam()}];
		}
	}
</script>

<style>
	h3 {
		color: #ee4035;
	}

	form {
		text-align: left;
		position: relative;
		display: flexbox;
		padding: 1em;
		max-width: 30em;
		margin: 0 auto;
		color: #eef2f3;
	}

	form *:not([type="radio"]) {
		width: 100%;
	}

	form input, form select, form option {
		background-color: #eef2f3;
		color: #242020;
		border: none;
		border-radius: 2px;
	}

	form > label {
		margin: 1em 0;
	}

	div {
		margin: 0.5em 0;
    }
    
    #status {
		color: #eef2f3;
        position: absolute;
        height: 1em;
        width: 100%;
        text-align: center;
        margin: auto;
        top: 0;
        bottom: 0;
    }

	/* display: grid allows stacking of all the explainers, so their height is always the maximum of their heights */
	#explainerContainer {
		display: grid;
		margin: 0;
	}

	#explainer {
		font-size: 0.8em;
		grid-column: 1;
		grid-row: 1;
		margin: 0;
	}
</style>

<main class="darkBackground">
	<form id="newQForm" on:submit|preventDefault={handleNewQuestion} class={newQFormVisible ? "visible" : "hidden"}>
		<div class="formItem">
			<label class="formLabel"><h3>New Question</h3></label>
			<input id="newQName" type="text" required placeholder="e.g., Which jam would you prefer?" bind:value={newQName} maxlength={MAX_FIELD_CHARS}>
		</div>
		<label>Options</label>
		{#each newQOptions as option}
			<div class="newQOption">
				<input bind:value={option.text} placeholder={option.placeholder} on:input={handleOptionUpdate} maxlength={MAX_FIELD_CHARS}>
			</div>
		{/each}
		<div class="radioSelect">
			{#each newQTypes as type}
				<input type="radio" bind:group={newQType} value={type} id={type.value} checked>
				<label class="clickable" for={type.value}><p>{type.text}</p></label>
			{/each}
		</div>
		<div id="explainerContainer">
			{#each newQTypes as type}
				<div id="explainer" class={type.value === newQType.value && newQFormVisible ? "visible" : "hidden"}>
					{@html type.explainer}
				</div>
			{/each}
		</div>
		<div>
			<button type=submit class="clickable">
				Create
			</button>
		</div>
	</form>
	<p id="status" class={newQFormVisible ? "hidden" : "visible"}>Submitting...</p>
</main>