/*
This file is (C) Captain ALM
Under the Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License
*/
const EntryData = []
function CreateEntry(name, videourl, videotype, start, end, duration) {
    return {
        name: name,
        videourl: videourl,
        videotype: videotype,
        start: Date.parse(start),
        end: Date.parse(end),
        duration : duration
    };
}
function CreateVideoPlaceholder(id) {
    let imgPH = document.createElement("img")
    imgPH.src = PlayImageURL
    imgPH.id = "play-"+id
    imgPH.alt = "Play Video"
    imgPH.width = 360
    imgPH.style.cursor = "pointer"
    if (document.addEventListener) {
        imgPH.addEventListener("click", function () {ActivateVideo(id);})
    } else {
        imgPH.setAttribute("onclick", "ActivateVideo("+id+");")
        imgPH.onclick = function () {ActivateVideo(id);}
    }
    document.getElementById("video-" + id).appendChild(imgPH)
}
function ActivateVideo(id) {
    let holder = document.getElementById("video-" + id)
    holder.removeChild(document.getElementById("play-"+id))
    let vid = document.createElement("video")
    vid.controls = true
    vid.width = 360
    let vids = document.createElement("source")
    vids.src = EntryData[id].videourl
    vids.type = EntryData[id].videotype
    let vida = document.createElement("a")
    vida.href = EntryData[id].videourl
    vida.innerText = "The Video"
    vid.appendChild(vids)
    vid.appendChild(vida)
    holder.appendChild(vid)
    if (vid.play) {
        vid.play()
    }
}
