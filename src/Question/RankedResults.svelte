<script>
    export let results = [];
    export let q;

    let maxVotes = 0;
    let winner = "";
    let winnerText = "";

    $: finalRound = results.reduce(function(a, b) {
        return a.round_num > b.round_num ? a : b;
    }, {"round_num": -1}).round_num;
    $: winner = results.filter(result => result.round_num == finalRound).reduce(function(a, b) {
        return a.votes > b.votes ? a : b;
    }, {"votes": -1});
    // TODO: this is not okay
    $: try {winnerText = q.options.filter(option => option.option_id == winner.option_id)[0].text} catch {winnerText = "loading..."};
</script>

<p>Runoff Winner: {winnerText}</p>
<p>(Visualization for these results is WIP)</p>