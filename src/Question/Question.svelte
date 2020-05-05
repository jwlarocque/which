<script>
    import ApprovalQuestion from "./ApprovalQuestion.svelte"
    import ApprovalResults from "./ApprovalResults.svelte"

    export let id;

    let q = {};
    let optionCounts = [];
    let votes;
    let ballots;

    let newVoteFormVisible = false;

    getQuestion(id);
    getBallots(id);

    // TODO: also retrieve user's current votes and fill them in
    async function getQuestion(question_id) {
        const res = await fetch("qs/q/" + question_id);
		const data = await res.json();

		if (res.ok) {
            // set votes first to avoid .fill race condition
            votes = Array(data.options.length);
            q = data;
            newVoteFormVisible = true;
            console.log(q);
		} else {
			throw new Error(data);
		}
    }

    async function getBallots(question_id) {
        const res = await fetch("qs/vs/" + question_id);
		const data = await res.json();

		if (res.ok) {
            ballots = data;
            updateResults(ballots);
		} else {
			throw new Error(data);
		}
    }

    function updateResults(ballots) {
        console.log(ballots);
        optionCounts = ballots.reduce(combineReducer, []).reduce(countReducer, []);
        console.log(optionCounts)
    }

    function combineReducer(vs, ballot) {
        return vs.concat(ballot.votes);
    }
    
    function countReducer(counts, vote) {
        if (counts[vote.option_id]) {
            counts[vote.option_id] += vote.state;
        } else {
            counts[vote.option_id] = vote.state;
        }
        return counts;
    }

    function voteStringFromState() {
        // TODO: this is not okay, but it does work
        // TODO: I think this might be too many anonymous functions in one line...
        //       it looks like Klingon...
        let ret = JSON.stringify(
            {"question_id": id, 
             "votes": votes.map( function(vote, index) { 
                 return {"option_id": index, "state": (vote === true ? 1 : (vote === false ? 0 : vote))}})})
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
    h3 {
		color: #ee4035;
	}
    
    form, #results {
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

<div class="darkBackground">
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
</div>
<div id="results">
    <h3>Results</h3>
    {#if q.type == "approval"}
        <ApprovalResults {q} {optionCounts}/>
    {:else if q.type == "runoff"}
        <p>runoff</p>
    {:else}
        <p>plurality</p>
    {/if}
</div>
