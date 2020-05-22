<script>
    export let results = [];
    export let q;

    let maxVotes = 0;
    let winner = "";

    $: finalRound = results.reduce(function(a, b) {
        return a.round_num > b.round_num ? a : b;
    }, {"round_num": -1}).round_num;
    $: winner = results.filter(result => result.round_num == finalRound).reduce(function(a, b) {
        return a.votes > b.votes ? a : b;
    }, {"votes": -1});
</script>

<p>Runoff Winner: {(q && q.options && q.options.length) ? q.options.filter(option => option.option_id == winner.option_id)[0].text : "loading..."}</p>
<p>(Visualization for these results is WIP)</p>