<script>
    import {qs} from "../stores.js"

    const MAX_OPTIONS = 50;
	const OPTION_PLACEHOLDERS = ["Apricot", "Rhubarb", "Sour Cherry", "Raspberry", "Hot Pepper", "Gooseberry", "Peach", "Quince", "Lingonberry", "Quince", "Cloudberry", "", "Strawberry", "Blackberry", "Blueberry", "Grape", "Orange Marmalade", "Plum", "Apple Butter", "Fig"]; // example options

	// this is the kind of function Gotham deserves
	function randomJam() {
		return "e.g., " + OPTION_PLACEHOLDERS[Math.floor(Math.random() * OPTION_PLACEHOLDERS.length)];
	}

    let newQFormVisible = true;
	let newQName = "";
	let newQOptions = [{id: 0, text: ""}, {id: 1, text: ""}, {id: 2, text: ""}];

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
        // TODO: send request to server
        newQFormVisible = false;
        if (validateQuestion()) {
            const res = await fetch("qs/new", {
                method: "POST",
                headers: {"Content-Type": "application/json",},
                body: JSON.stringify({"name": newQName, "options": newQOptions}),
            });
            const data = await res.json();

            if (res.ok) {
                clearNewQForm();
                newQFormVisible = true;
                qs.update(value => [...value, data]);
            } else {
                newQFormVisible = true;
                throw new Error(data);
            }
        } else {
            newQFormVisible = true;
        }
	}

	function handleOptionUpdate() {
		if (newQOptions[newQOptions.length - 1].text !== "" && newQOptions.length < MAX_OPTIONS) {
			newQOptions = [...newQOptions, {id: newQOptions.length, text: "", placeholder: randomJam()}];
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
		max-width: 30em;
		margin: 0 auto;
	}

	form * {
		width: 100%;
	}

	#formBackground {
		width: 100%;
        background-color: #242020;
        position: relative;
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
</style>

<main>
    <div id="formBackground">
		<form id="newQForm" on:submit|preventDefault={handleNewQuestion} class={newQFormVisible ? "visible" : "hidden"}>
			<div class="formItem">
				<label class="formLabel"><h3>New Question</h3></label>
				<input id="newQName" type="text" required placeholder="e.g., Which jam would you prefer?" bind:value={newQName}>
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
        <p id="status" class={newQFormVisible ? "hidden" : "visible"}>Submitting...</p>
	</div>
</main>