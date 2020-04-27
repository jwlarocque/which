<script>
    import ApprovalQuestion from "./ApprovalQuestion.svelte"

    export let id;

    let q = {};
    let votes;

    let newVoteFormVisible = false;

    getQuestion(id);

    async function getQuestion(question_id) {
        const res = await fetch("qs/q/" + question_id);
		const data = await res.json();

		if (res.ok) {
            // set votes first to avoid .fill race condition
            votes = Array(data.options.length);
            q = data;
            newVoteFormVisible = true;
		} else {
			throw new Error(data);
		}
    }

    function voteStringFromState() {
        // converts votes = [true, false, true true] 
        //     and produces [0,           2,   3]
        // and puts it in a JSON string along with the question_id
        // TODO: I think this might be too many anonymous functions in one line...
        //       it looks like Klingon...
        let ret = JSON.stringify(
            {"question_id": id, 
             "votes": votes.map( function(vote, index) { 
                 return {"id": index, "state": vote}})})
        console.log(ret)
        return ret
    }

    async function handleNewVote() {
        newVoteFormVisible = false;
        const res = await fetch("qs/vote", {
            method: "POST",
            headers: {"Content-Type": "application/json",},
            body: voteStringFromState(),
        });
        const data = await res.json();

        if (res.ok) {
            newVoteFormVisible = true;
            console.log(res) // TODO: remove debug
            // TODO: display/update vote results
        } else {
            newVoteFormVisible = true;
            throw new Error(data.message); // TODO: improve and replicate this error pattern
        }
    }
</script>

<style>
    main :global(*) {
        /*background-color: red;*/
    }

    h3 {
		color: #ee4035;
	}
    
    form {
        text-align: left;
        max-width: 30em;
        margin: 0 auto;
        padding: 1em;
    }

    button {
        width: 100%;
        margin: 1.4em 0 1em;
    }

    .statusMessage {
        text-align: center;
    }
</style>

<main class="darkBackground">
    <form id="newVoteForm" on:submit|preventDefault={handleNewVote}>
        {#if q.name}
            <h3>{q.name}</h3>
            {#if q.type == "approval"}
                <ApprovalQuestion {q} {votes}/>
            {:else if q.type == "runoff"}
                <p>runoff</p>
            {:else}
                <p>plurality</p>
            {/if}
        {:else}
            <p class="statusMessage">Loading...</p>
        {/if}
        <button type=submit class="clickable">
            Submit
        </button>
    </form>
</main>