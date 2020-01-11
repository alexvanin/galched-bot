// queue is used to play videos
var queue
// player is a iframe youtube player
var player
// vid is a current playing video id
var vid = ""

function updateTable() {
    var myTableDiv = document.getElementById("video_queue")
    myTableDiv.innerHTML = ""

    var table = document.createElement('TABLE')
    var tableBody = document.createElement('TBODY')
    table.border = '1'
    table.appendChild(tableBody);
    
    for (i = 0; i < queue.length; i++) {
        var tr = document.createElement('TR');

        var td = document.createElement('TD');
        td.appendChild(document.createTextNode(queue[i].Duration));
        tr.appendChild(td);

        var td = document.createElement('TD');
        td.appendChild(document.createTextNode(queue[i].Title));
        tr.appendChild(td)

        var td = document.createElement('TD');
        td.appendChild(document.createTextNode( queue[i].From));
        tr.appendChild(td)

        var td = document.createElement('TD');
        td.appendChild(document.createTextNode( queue[i].Views));
        tr.appendChild(td)

        var td = document.createElement('TD');
        td.appendChild(document.createTextNode(((queue[i].Upvotes/(queue[i].Upvotes + queue[i].Downvotes)) * 100).toFixed(2) + '%'));
        tr.appendChild(td)

        tableBody.appendChild(tr);
    }

    myTableDiv.appendChild(table)
}

function bgUpdateQueue() {
    initQueue()
    updateTable()

    if (vid === "" && queue.length > 0) {
        vid = queue[0].ID
        player.loadVideoById(vid)
    }
}

function initQueue() {
    var request = new XMLHttpRequest()
    request.open('GET', './queue', false)
    request.onload = function () {
        queue = JSON.parse(this.response)
    }
    request.send()
    updateTable()
}

function updateQueue() {
    var request = new XMLHttpRequest()
    request.open('POST', './queue', false)
    request.onload = function () {
        queue = JSON.parse(this.response)
    }
    request.send(vid)
    updateTable()
}

function nextVideo() {
    updateQueue()
    if (queue.length > 0) {
        vid = queue[0].ID
    } else {
        vid = ""
    }
        
    player.loadVideoById(vid)
    updateTable()
}

function onYouTubeIframeAPIReady() {
    if (queue.length != 0) {
        vid = queue[0].ID
    }
    console.log(vid)
    player = new YT.Player('player', {
        height: '390',
        width: '640',
        videoId: vid,
        events: {
            'onReady': onPlayerReady,
            'onStateChange': onPlayerStateChange
        }
    });

}

function onPlayerReady(event) {
    event.target.playVideo();
}

function onPlayerStateChange(event) {
    if (event.data === 0) {
        nextVideo()
    }
}

initQueue()
setInterval(bgUpdateQueue, 4000);
var tag = document.createElement('script');
tag.src = "https://www.youtube.com/iframe_api";
var firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
