<script>
    import ApprovalQuestion, {approvalVotesFromBallot} from "./ApprovalQuestion.svelte"
    import SingleRoundResults from "./SingleRoundResults.svelte"
    import RankedQuestion, {rankedVotesFromBallot} from "./RankedQuestion.svelte"
    import RankedResults from "./RankedResults.svelte"
    import PluralityQuestion, {pluralityVotesFromBallot} from "./PluralityQuestion.svelte"

    export let id;

    let q = {};
    let optionCounts = [];
    let votes = [];
    let results;

    let newVoteFormVisible = false;

    getQuestion(id);
    getResults(id);

    // TODO: also retrieve user's current votes and fill them in
    async function getQuestion(question_id) {
        const res = await fetch("qs/q/" + question_id);
		const data = await res.json();

		if (res.ok) {
            q = data;
            getBallot(id);
		} else {
			throw new Error(data);
		}
    }

    // get the user's ballot, if one exists
    async function getBallot(question_id) {
        const res = await fetch("qs/b/" + question_id);
        let ballot;
        if (res.ok) {
            try {
                ballot = await res.json();
            } catch {
                ballot = {"votes": []};
            }
        } else {
            // if bad/no response, we can just make the user
            // fill in the form from scratch
            ballot = {"votes": []};
        }

        if (q.type == 0) {
            votes = approvalVotesFromBallot(q, ballot);
        } else if (q.type == 1) {
            votes = rankedVotesFromBallot(q, ballot);
        } else if (q.type == 2) {
            votes = pluralityVotesFromBallot(q, ballot);
        } else {
            votes = [];
        }
        newVoteFormVisible = true;
    }

    async function getResults(question_id) {
        const res = await fetch("qs/rs/" + question_id);
		const data = await res.json();

		if (res.ok) {
            results = data;
		} else {
			throw new Error(data);
		}
    }

    function voteStringFromState() {
        // TODO: this is not okay, but it does work
        // TODO: I think this might be too many anonymous functions in one line...
        //       it looks like Lisp...
        let ret = JSON.stringify(
            {"question_id": id, 
             "votes": votes.map( function(vote, index) { 
                 return {"option_id": index, "state": (vote === true ? 1 : (vote === false ? 0 : vote))}})})
        return ret
    }

    async function handleNewVote() {
        newVoteFormVisible = false;
        const res = await fetch("qs/vote", {
            method: "POST",
            headers: {"Content-Type": "application/json",},
            body: voteStringFromState(),
        });

        if (res.ok) {
            newVoteFormVisible = true;
            getResults(id);
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

    form button {
        width: 100%;
        margin: 1.4em 0 1em;
    }

    .statusMessage {
        text-align: center;
    }

    #results div {
        width: 100%;
        display: inline-flex;
        justify-content: space-between;
    }

    #results button {
        max-height: 2em;
        margin: auto 0;
        cursor: pointer;
    }
</style>

<div class="darkBackground">
    <form id="newVoteForm" on:submit|preventDefault={handleNewVote}>
        {#if q.name}
            <h3>{q.name}</h3>
            {#if q.type == 0} <!-- TODO: question type enum -->
                <ApprovalQuestion {q} {votes}/>
            {:else if q.type == 1}
                <RankedQuestion bind:q={q} bind:votes={votes}/>
            {:else if q.type == 2}
                <PluralityQuestion {q} {votes}/>
            {:else}
                <p>error</p>
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
    <div>
        <h3>Results</h3>
        <button class="button" on:click={getResults(id)}>Refresh</button>
    </div>
    {#if q.type == 0}
        <SingleRoundResults {q} {results}/>
    {:else if q.type == 1}
        <RankedResults {q} {results}/>
    {:else if q.type == 2}
         <SingleRoundResults {q} {results}/>
    {:else}
        <p>error</p>
    {/if}
</div>
