/*
This file is (C) Captain ALM
Under the Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License
*/
const EntryData = []
SetupJSTheme()
function CreateEntry(id, name, videourl, videotype, start, end, duration) {
    EntryData[id] = {
        name: name,
        videourl: videourl,
        videotype: videotype,
        start: Date.parse(start),
        end: Date.parse(end),
        duration : parseInt(duration, 10)
    };
}
function CreateVideoPlaceholder(id) {
    let imgPH = document.createElement("img")
    imgPH.src = PlayImageURL
    imgPH.id = "play-"+id
    imgPH.alt = "Play Video"
    imgPH.title = "Play"
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
    if (vid.play) {vid.play();}
}
function SetupJSTheme() {
    let th = document.getElementById("theme")
    th.href = "#"
    if (document.addEventListener) {
        th.addEventListener("click", ToggleTheme)
    } else {
        th.setAttribute("onclick", "ToggleTheme();")
        th.onclick = ToggleTheme
    }
}
function ToggleTheme() {
    let th = document.getElementById("theme")
    let thimg = document.getElementById("theme-img")
    let thsty = document.getElementById("style-theme")
    let logo = document.getElementById("logo")
    if (document.getElementById("so-theme")) {
        thimg.src = SunImageURL
        thimg.alt = "()"
        th.title = "Switch to Light Mode"
        document.getElementById("so-form").removeChild(document.getElementById("so-theme"))
        logo.href = "?"
        thsty.src = CssDarkURL
    } else {
        thimg.src = MoonImageURL
        thimg.alt = "{"
        th.title = "Switch to Dark Mode"
        let thi = document.createElement("input")
        thi.name = "light"
        thi.type = "hidden"
        thi.id = "so-theme"
        document.getElementById("so-form").appendChild(thi)
        logo.href = "?light"
        thsty.src = CssLightURL
    }
}
