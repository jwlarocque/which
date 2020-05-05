<!-- Loosely based on 
         https://dev.to/richharris/a-new-technique-for-making-responsive-javascript-free-charts-gmp 
         and https://pancake-charts.surge.sh/ by Rich Harris-->
<script>
    export let data;
    export let xLabels = [];
    export let yLabels = [];

    $: yMax = Math.max(...data.map(datum => datum.x));
    $: updateXLabels(yMax);

    function updateXLabels(max) {
        xLabels = [...Array(max ? max + 1 : 10).keys()]
    }

    const barHeight = 2; // em
    const barMinWidth = 1; // %
    const barMaxWidth = 60; // %
    const horizGap = 0.5; // em
</script>

<style>
    .barChart {
        width: 100%;
        max-width: 30em;
        margin: 1em auto;
        padding: 0 1em;
        position: relative;
    }

    svg * {
		vector-effect: non-scaling-stroke;
	}

    .bar {
        fill: #242020;
    }

    .yLabels {
        position: absolute;
        top: 0;
        height: 100%;
        width: 100%;
    }

    .label {
        position: absolute;
        margin: 0;
        padding: 0;
    }

    .xLabels {
        height: 0;
        margin: 0;
        padding: 0;
        top: -1em;
        display: inline-flex;
        justify-content: space-between;
        position: relative;
    }

    .xLabels p {
        position: absolute;
        margin: 0;
    }
</style>

<div class="barChart">
    <div style="postion:relative;">
        <svg width="100%" height="{data.length * (barHeight + horizGap) - horizGap}em" preserveaspectratio="none">
            {#each data as datum, i}
                <rect 
                    class="bar" 
                    width="{(Math.max(barMinWidth, barMaxWidth * data[i].x / yMax))}%"
                    height="{barHeight}em"
                    x="{100 - barMaxWidth}%"
                    y="{(barHeight * i) + (horizGap * (i))}em"
                />
            {/each}
        </svg>
        <div class="yLabels" width="{100 - barMaxWidth}%">
            {#each yLabels as label, i}
                <p 
                    class="label"
                    width="{100 - barMaxWidth}%"
                    style="top: calc({i * 100 / yLabels.length}% + {((i / yLabels.length) * (horizGap * 0.5 - barHeight) * 0.5) + 0.5}em)"
                >{label}</p>
            {/each}
        </div>
    </div>
    <div class="xLabels" style="width:{barMaxWidth}%; margin-left:{100 - barMaxWidth}%">
        {#each xLabels as label, i}
            <p style="left: calc({i * 100 / (xLabels.length - 1)}% - 5px)">{label}</p>
        {/each}
    </div>
</div>
